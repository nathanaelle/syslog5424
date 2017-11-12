package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
)

type (
	Listener interface {
		net.Listener
	}

	Collector struct {
		// length of the queue to the receiver queue
		QueueLen int

		scan     *bufio.Scanner
		pipeline chan []byte
	}

	Receiver struct {
		listener  Listener
		transport Transport
		pipeline  chan messageErrorPair
		end       chan struct{}
	}

	messageErrorPair struct {
		m MessageImmutable
		e error
	}
)

func Collect(network, address string) (*Receiver, error) {
	return (Collector{
		QueueLen: 100,
	}).Collect(network, address, nil)
}

func (d Collector) Collect(network, address string, t Transport) (*Receiver, error) {
	var pipeline chan messageErrorPair
	var c Listener
	var err error

	switch network {
	case "unix":
		if t == nil {
			t = T_ZEROENDED
		}
		c, err = unix_coll(network, address)

	case "unixgram":
		if t == nil {
			t = T_ZEROENDED
		}
		c, err = unixgram_coll(network, address)

	case "tcp", "tcp6", "tcp4":
		if t == nil {
			t = T_LFENDED
		}
		c, err = tcp_coll(network, address)

	default:
		return nil, errors.New("unknown network for Collector : " + network)
	}

	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errors.New("No Connection established")
	}

	switch d.QueueLen <= 0 {
	case true:
		pipeline = make(chan messageErrorPair)

	case false:
		pipeline = make(chan messageErrorPair, d.QueueLen)
	}

	return NewReceiver(c, pipeline, t), nil
}

func NewReceiver(listener Listener, pipeline chan messageErrorPair, t Transport) *Receiver {
	r := &Receiver{
		listener:  listener,
		pipeline:  pipeline,
		transport: t,
		end:       make(chan struct{}),
	}

	go r.run_queue()

	return r
}

func (r *Receiver) run_queue() {
	defer r.listener.Close()
	defer close(r.pipeline)

	for {
		select {
		case <-r.end:
			return

		default:
			conn, err := r.listener.Accept()
			if err != nil {
				panic(err)
			}

			go r.tokenize(conn)
		}
	}

}

func (r *Receiver) tokenize(conn io.ReadCloser) {
	defer conn.Close()

	var buffer []byte
	var data []byte
	var eof bool

	for {
		if buffer == nil {
			buffer = make([]byte, 1<<20)
		}

		size, err := conn.Read(buffer)
		if err == io.EOF {
			eof = true
			err = nil
		}
		if err != nil {
			panic(err)
		}

		data, buffer = buffer[0:size], buffer[size:]
		loop := true
		for loop {
			msg, rest, m_err := Parse(data, r.transport, eof)

			log.Printf("L {%q} {%q} %v", msg, rest, m_err)

			if len(rest) == 0 {
				loop = false
			}
			if err != nil && rest == nil {
				break
			}
			data = rest

			r.pipeline <- messageErrorPair{msg, m_err}
		}

		if eof == true {
			return
		}
	}
}

func (r *Receiver) Receive() (MessageImmutable, error, bool) {
	pair, end := <-r.pipeline

	return pair.m, pair.e, end
}

// terminate the log_collector goroutine
func (r *Receiver) End() {
	close(r.end)
}
