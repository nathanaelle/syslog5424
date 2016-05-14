package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"time"
	"errors"
)

type	(
	Dialer struct {
		// length of the queue to the log_sender goroutine
		QueueLen	int

		// delay to flush the queue
		FlushDelay	time.Duration
	}
)


// Dial opens a connection to the syslog daemon
// network can be "stdio", "unix", "unixgram", "tcp", "tcp4", "tcp6"
// used Transport is the "common" transport for the network.
// QueueLen is preset to 100 Message
// FlushDelay is preset to 500ms
func Dial(network, address string) (*Sender,error) {
	return (Dialer{
		QueueLen:	100,
		FlushDelay:	500*time.Millisecond,
	}).Dial(network, address, nil)
}


// Dial opens a connection to the syslog daemon
// network can be "stdio", "unix", "unixgram", "tcp", "tcp4", "tcp6"
// Transport can be nil.
// if Transport is nil the "common" transport for the wished network is used.
func (d Dialer) Dial(network, address string, t Transport) (*Sender,error) {
	var pipeline	chan Message
	var ticker	<-chan time.Time
	var c		Conn

	switch network {
	case "stdio":
		if t == nil {
			t = new(T_LFENDED)
		}
		c = stdio_dial(address)

	case "local":
		if t == nil {
			t = new(T_ZEROENDED)
		}
		c = local_dial("", address)

	case "unix", "unixgram":
		if t == nil {
			t = new(T_ZEROENDED)
		}
		c = local_dial(network, address)

	case "tcp", "tcp6", "tcp4":
		if t == nil {
			t = new(T_LFENDED)
		}
		c = tcp_dial(network, address)

	default:
		return nil, errors.New("unknown network for Dial : "+network)
	}

	if c == nil {
		return nil, errors.New("No Connection established")
	}

	switch d.QueueLen <= 0 {
	case true:
		pipeline = make(chan Message)

	case false:
		pipeline = make(chan Message, d.QueueLen)
	}

	switch {
	case d.FlushDelay <= time.Millisecond:
		// less than 1ms => disable auto flush
		ticker	= make(chan time.Time)

	default:
		ticker	= time.Tick(d.FlushDelay)
	}

	t.SetConn(c)

	return NewSender(t, pipeline, ticker), nil
}
