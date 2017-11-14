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
		connector     Connector
		output        io.WriteCloser
		end_asked     chan struct{}
		end_completed chan struct{}
		ticker        <-chan time.Time
		transport     Transport
		err_chan      chan error
		lock          *sync.Mutex
		queue         *net.Buffers
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
		end_asked:     make(chan struct{}),
		end_completed: make(chan struct{}),
		connector:     output,
		ticker:        ticker,
		transport:     transport,
		err_chan:      make(chan error, 1),
		lock:          new(sync.Mutex),
		queue:         new(net.Buffers),
	}
	*s.queue = make([][]byte, 0, 1000)

	go s.run_queue()

	return s, s.err_chan
}

func (c *Sender) flush_queue() {
	if c.output == nil {
		var err error
		if c.output, err = c.connector.Connect(); err != nil {
			c.err_chan <- err
			return
		}
	}

	//log.Printf("<--\tget lock for %d items", len(*c.queue))
	c.lock.Lock()
	defer c.lock.Unlock()

	for len(*c.queue) > 0 {
		_, err := c.queue.WriteTo(c.output)

		switch t_err := err.(type) {
		case nil:
		case *net.OpError:
			if s_err, ok := t_err.Err.(*os.SyscallError); ok && s_err.Err == syscall.ENOBUFS {
				return
			}
			c.err_chan <- err
			c.output = nil
			return

		default:
			c.err_chan <- err
			c.output = nil
			return
		}
	}

	if len(*c.queue) == 0 && cap(*c.queue) == 0 {
		*c.queue = make([][]byte, 0, 100)
	}
	//log.Printf("<--\tget unlock")
}

func (c *Sender) run_queue() {
	defer func() {
		for len(*c.queue) > 0 {
			c.flush_queue()
		}
		c.output.Close()
		close(c.end_completed)
		close(c.err_chan)
	}()

	for {
		select {
		case <-c.ticker:
			c.flush_queue()

		case _, opened := <-c.end_asked:
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
	close(c.end_asked)
	<-c.end_completed
}

func (a *Addr) String() string {
	return a.network + "!" + a.address
}

func (a *Addr) Network() string {
	return a.network
}
