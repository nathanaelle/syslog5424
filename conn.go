package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"net"
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
		output        WriteCloser
		pipeline      chan []byte
		end_completed chan struct{}
		ticker        <-chan time.Time
		transport     Transport
		err_chan      chan error
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
func NewSender(output Connector, transport Transport, pipeline chan []byte, ticker <-chan time.Time) (*Sender, <-chan error) {
	s := &Sender{
		pipeline:      pipeline,
		end_completed: make(chan struct{}),
		connector:     output,
		ticker:        ticker,
		transport:     transport,
		err_chan:      make(chan error, 1),
	}

	go s.run_queue()

	return s, s.err_chan
}

func (c *Sender) run_queue() {
	queue := new(net.Buffers)
	*queue = make([][]byte, 0, 1000)

	defer func() {
		c.output.Close()
		close(c.end_completed)
		close(c.err_chan)
	}()

	for {
		select {
		case <-c.ticker:
			if c.output == nil {
				var err error
				if c.output, err = c.connector.Connect(); err != nil {
					c.err_chan <- err
					continue
				}
			}

			for len(*queue) > 0 {
				if _, err := queue.WriteTo(c.output); err != nil {
					c.err_chan <- err
					c.output = nil
					break
				}
			}
			if len(*queue) == 0 && cap(*queue) == 0 {
				*queue = make([][]byte, 0, 1000)
			}

		case msg, opened := <-c.pipeline:
			if !opened && msg == nil {
				switch len(*queue) {
				case 0:
					return
				default:
					continue
				}
			}

			*queue = append(*queue, msg)
		}
	}
}

// send a Message to the log_sender goroutine
func (c *Sender) Send(m Message) (err error) {
	msg, err := m.Marshal5424()
	c.pipeline <- c.transport.Encode(msg)
	return
}

// terminate the log_sender goroutine
func (c *Sender) End() {
	close(c.pipeline)
	<-c.end_completed
}

func (a *Addr) String() string {
	return a.network + "!" + a.address
}

func (a *Addr) Network() string {
	return a.network
}
