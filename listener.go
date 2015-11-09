package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
	"bufio"
	"errors"
	"crypto/tls"
)


type	(



	Listener interface {
		net.Listener
	}


	Collector struct {
		// length of the queue to the receiver queue
		QueueLen	int

		scan		*bufio.Scanner
		pipeline	chan []byte

		TLSConf		tls.Config
	}

	Receiver interface {
		SetTransport(Transport)
		RunQueue(chan []byte)

		Receive() ([]byte, bool)

		// terminate the log_collector goroutine
		End()
	}

)



func (d Collector)Collect(network,address string, t Transport) (Receiver,error) {
	var pipeline chan []byte
	var c Receiver

	switch network {
	case "unix":
		if t == nil {
			t = new(T_ZEROENDED)
		}
		c = unix_coll(network, address)

	case "unixgram":
		if t == nil {
			t = new(T_ZEROENDED)
		}
		c = unixgram_coll(network, address)

	case "tcp", "tcp6", "tcp4":
		if t == nil {
			t = new(T_LFENDED)
		}
		c = tcp_coll(network, address)

	/*
	case "tls", "tls6", "tls4":
		if t == nil {
			t = new(T_LFENDED)
		}
		c = tls_coll(network, address)
	*/

	default:
		return nil, errors.New("unknown network for Collector : "+network)
	}

	if c == nil {
		return nil, errors.New("No Connection established")
	}

	switch d.QueueLen <= 0 {
	case true:
		pipeline = make(chan []byte)

	case false:
		pipeline = make(chan []byte, d.QueueLen)
	}

	c.SetTransport(t)
	go c.RunQueue(pipeline)

	return c,nil
}
