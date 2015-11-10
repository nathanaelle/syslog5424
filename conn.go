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
		output		Conn
		pipeline	chan Message
		ticker		<-chan time.Time
	}


	Addr	struct {
		network	string
		address string
	}
)

// Create a new sender
func NewSender(output Conn, pipeline chan Message, ticker <-chan time.Time) (*Sender) {
	s := &Sender {
		pipeline:	pipeline,
		output:		output,
		ticker:		ticker,
	}

	go s.run_queue()

	return s
}



func (c *Sender) run_queue() {
	if err := c.output.Redial(); err != nil {
		log.Fatal(err)
	}

	defer c.output.Close()

	for {
		select {
		case <-c.ticker:
			c.output.Flush()

		case msg, opened := <-c.pipeline:
			if !opened {
				return
			}

			_, err := c.output.Write(msg.Marshal5424())
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
}


func (a *Addr) String() string {
	return a.network + "!" + a.address
}

func (a *Addr) Network() string {
	return a.network
}
