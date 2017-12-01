package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

type (
	Connector interface {
		Connect() (WriteCloser, error)
	}

	ConnectorFunc func() (WriteCloser, error)

	// generic interface describing a Connection
	WriteCloser interface {
		io.Writer
		io.Closer
	}

	// Sender describe the generic algorithm for sending Message through a connection
	Sender struct {
		connector    Connector
		output       io.WriteCloser
		endAsked     chan struct{}
		endCompleted chan struct{}
		ticker       <-chan time.Time
		transport    Transport
		errChan      chan error
		lock         *sync.Mutex
		queue        *net.Buffers
	}

	Addr struct {
		network string
		address string
	}
)

func (f ConnectorFunc) Connect() (WriteCloser, error) {
	return f()
}

// Create a new sender
func NewSender(output Connector, transport Transport, ticker <-chan time.Time) (*Sender, <-chan error) {
	s := &Sender{
		endAsked:     make(chan struct{}),
		endCompleted: make(chan struct{}),
		connector:    output,
		ticker:       ticker,
		transport:    transport,
		errChan:      make(chan error, 1),
		lock:         new(sync.Mutex),
		queue:        new(net.Buffers),
	}
	*s.queue = make([][]byte, 0, 1000)

	go s.runQueue()

	return s, s.errChan
}

func (c *Sender) flushQueue() {
	if c.output == nil {
		var err error
		if c.output, err = c.connector.Connect(); err != nil {
			c.errChan <- err
			return
		}
	}

	//log.Printf("<--\tget lock for %d items", len(*c.queue))
	c.lock.Lock()
	defer c.lock.Unlock()

	for len(*c.queue) > 0 {
		_, err := c.queue.WriteTo(c.output)

		switch typeErr := err.(type) {
		case nil:
		case *net.OpError:
			if sysErr, ok := typeErr.Err.(*os.SyscallError); ok && sysErr.Err == syscall.ENOBUFS {
				return
			}
			c.errChan <- err
			c.output = nil
			return

		default:
			c.errChan <- err
			c.output = nil
			return
		}
	}

	if len(*c.queue) == 0 && cap(*c.queue) == 0 {
		*c.queue = make([][]byte, 0, 100)
	}
	//log.Printf("<--\tget unlock")
}

func (c *Sender) runQueue() {
	defer func() {
		for len(*c.queue) > 0 {
			c.flushQueue()
		}
		c.output.Close()
		close(c.endCompleted)
		close(c.errChan)
	}()

	for {
		select {
		case <-c.ticker:
			c.flushQueue()

		case _, opened := <-c.endAsked:
			if !opened {
				return
			}
		}
	}
}

// send a Message to the log_sender goroutine
func (c *Sender) Send(m Message) (err error) {
	var msg []byte

	msg, err = m.Marshal5424()
	if err != nil {
		return err
	}

	c.lock.Lock()
	*c.queue = append(*c.queue, c.transport.Encode(msg))
	c.lock.Unlock()

	return
}

// terminate the log_sender goroutine
func (c *Sender) End() {
	close(c.endAsked)
	<-c.endCompleted
}

func (a *Addr) String() string {
	return a.network + "!" + a.address
}

func (a *Addr) Network() string {
	return a.network
}
