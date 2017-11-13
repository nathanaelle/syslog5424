package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"io"
	"net"
//	"log"
)

type (
	Listener interface {
		// set a deadline for Accept()
		SetDeadline(t time.Time) error

		// Accept waits for and returns the next DataReader to the listener.
	         Accept() (DataReader, error)

	         // Close closes the listener.
	         Close() error
	}

	DataReader interface {
		io.Reader
		io.Closer

		// RemoteAddr returns the remote network address.
		RemoteAddr() net.Addr
	}


	Collector struct {
		// length of the queue to the receiver queue
		QueueLen int

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

const	readBuffer = 1<<18

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
		return nil, ErrorNoConnecion
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

	var eof bool

	done := 0
	count := 0
	total := 0
	buffer := make([]byte, readBuffer)
	for {
		read_len, err := conn.Read(buffer[done:])
		total += read_len
		if err == io.EOF {
			eof = true
			err = nil
		}
		//log.Printf("EOF\t%v %v %v %v", count, total, read_len, err)

		if err != nil {
			panic(err)
		}
		if read_len == 0 && eof {
			return
		}
		if read_len == 0 {
			continue
		}

		read_len += done
		data := buffer[0:read_len]
		for {
			msg, rest, m_err := Parse(data, r.transport, eof)

			if rest == nil {
				//log.Printf("NIL\t{%q} {%q} %v %v %v", msg, data[0:10], len(rest), rest == nil, m_err)
				//log.Printf("NIL\t%v %v %v %v", count, total, read_len, m_err)
				break
			}
			data = rest

			count++
			r.pipeline <- messageErrorPair{msg, m_err}
			if len(rest) == 0 {
				break
			}
		}

		done = len(data)
		buffer = buffer[read_len-done:]

		if len(buffer) < 500 {
			old := buffer
			buffer = make([]byte, readBuffer)
			copy(buffer[0:len(old)], old)
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
