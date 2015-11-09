package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"log"
	"time"
)


type (

	Conn interface {
		io.Reader
		io.Writer
		io.Closer

		// reconnect a lost connection
		// Redial MUST wait in case of temporary errors
		// Redial MUST return in case of permanent error
		Redial() error

		Flush() error
	}


	Sender struct {
		output		Conn
		pipeline	chan Message
		ticker		<-chan time.Time
	}


)


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
