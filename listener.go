package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"net"
	//	"log"
)

type (
	Listener interface {
		// set a deadline for Accept()
		// SetDeadline(t time.Time) error

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

const readBuffer = 1 << 18

// if Transport is nil then the function returns nil, nil
// this case may occurs when transport is unknown at compile time
//
// the returned `<-chan error` is used to collect errors than may occur in goroutine
func NewReceiver(listener Listener, queue_len int, t Transport) (*Receiver, <-chan error) {
	var pipeline chan messageErrorPair

	if t == nil {
		return	nil, nil
	}

	if queue_len <= 0 {
		pipeline = make(chan messageErrorPair)
	} else {
		pipeline = make(chan messageErrorPair, queue_len)
	}

	r := &Receiver{
		listener:  listener,
		pipeline:  pipeline,
		transport: t,
		end:       make(chan struct{}),
	}

	chan_err := make(chan error, 10)

	go r.run_queue(chan_err)

	return r, chan_err
}

func (r *Receiver) run_queue(chan_err chan<- error) {
	defer r.listener.Close()
	defer close(r.pipeline)

	for {
		select {
		case <-r.end:
			return

		default:
			conn, err := r.listener.Accept()
			if err != nil {
				chan_err <- err
			}

			go r.tokenize(conn, chan_err)
		}
	}

}

func (r *Receiver) tokenize(conn io.ReadCloser, chan_err chan<- error) {
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
			chan_err <- err
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

// Read an incoming syslog message and a possible error that occured during the decoding of this syslog message
func (r *Receiver) Receive() (MessageImmutable, error, bool) {
	pair, end := <-r.pipeline

	return pair.m, pair.e, end
}

// terminate the log_collector goroutine
func (r *Receiver) End() {
	close(r.end)
}
