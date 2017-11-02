package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"log"
	"time"
)

type (
	// generic interface describing a Connection
	Conn interface {
		io.Reader
		io.Writer
		io.Closer

		// reconnect a lost connection
		// Redial MUST wait in case of temporary errors
		// Redial MUST return in case of permanent error
		Redial() error

		// flush all the remaining buffers
		Flush() error
	}

	// Sender describe the generic algorithm for sending Message through a connection
	Sender struct {
		output        Conn
		pipeline      chan Message
		end_completed chan struct{}
		ticker        <-chan time.Time
	}

	Addr struct {
		network string
		address string
	}
)

// Create a new sender
func NewSender(output Conn, pipeline chan Message, ticker <-chan time.Time) *Sender {
	s := &Sender{
		pipeline:      pipeline,
		end_completed: make(chan struct{}),
		output:        output,
		ticker:        ticker,
	}

	go s.run_queue()

	return s
}

func (c *Sender) run_queue() {
	defer func() {
		c.output.Close()
		close(c.end_completed)
	}()

	if err := c.output.Redial(); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-c.ticker:
			c.output.Flush()

		case msg, opened := <-c.pipeline:
			if !opened {
				return
			}

			raw, err := msg.Marshal5424()
			if err != nil {
				log.Fatal(err)
			}
			_, err = c.output.Write(raw)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// send a Message to the log_sender goroutine
func (c *Sender) Send(m Message) {
	c.pipeline <- m
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
