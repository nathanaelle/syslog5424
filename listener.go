package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"io"
	"net"
	//	"log"
)

type (
	// Listener decribe a generic way to Listen for an incoming connexion
	Listener interface {
		// set a deadline for Accept()
		// SetDeadline(t time.Time) error

		// Accept waits for and returns the next DataReader to the listener.
		Accept() (DataReader, error)

		// Close closes the listener.
		Close() error
	}

	// DataReader describe an incoming connexion
	DataReader interface {
		io.Reader
		io.Closer

		// RemoteAddr returns the remote network address.
		RemoteAddr() net.Addr
	}

	// Receiver describe how message are received and decoded
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

const readBuffer = 1 << 16

// NewReceiver create a new Listener
//
// if Transport is nil then the function returns nil, nil
// this case may occurs when transport is unknown at compile time
//
// the returned `<-chan error` is used to collect errors than may occur in goroutine
func NewReceiver(listener Listener, lenQueue int, t Transport) (*Receiver, <-chan error) {
	var pipeline chan messageErrorPair

	if t == nil {
		return nil, nil
	}

	if lenQueue <= 0 {
		pipeline = make(chan messageErrorPair)
	} else {
		pipeline = make(chan messageErrorPair, lenQueue)
	}

	r := &Receiver{
		listener:  listener,
		pipeline:  pipeline,
		transport: t,
		end:       make(chan struct{}),
	}

	chanErr := make(chan error, 10)

	go r.runQueue(chanErr)

	return r, chanErr
}

func (r *Receiver) runQueue(chanErr chan<- error) {
	defer r.listener.Close()
	defer close(r.pipeline)

	for {
		select {
		case <-r.end:
			return

		default:
			conn, err := r.listener.Accept()
			if err != nil {
				chanErr <- err
			}

			go r.tokenize(conn, chanErr)
		}
	}

}

func (r *Receiver) tokenize(conn io.ReadCloser, chanErr chan<- error) {
	defer conn.Close()

	var eof bool

	done := 0
	count := 0
	total := 0
	buffer := make([]byte, readBuffer)
	for {
		lenRead, err := conn.Read(buffer[done:])
		total += lenRead
		if err == io.EOF {
			eof = true
			err = nil
		}
		//log.Printf("EOF\t%v %v %v %v", count, total, lenRead, err)

		if err != nil {
			chanErr <- err
		}
		if lenRead == 0 && eof {
			return
		}
		if lenRead == 0 {
			continue
		}

		lenRead += done
		data := buffer[0:lenRead]
		for {
			msg, rest, parseErr := Parse(data, r.transport, eof)

			if rest == nil {
				//log.Printf("NIL\t{%q} {%q} %v %v %v", msg, data[0:10], len(rest), rest == nil, parseErr)
				//log.Printf("NIL\t%v %v %v %v", count, total, lenRead, parseErr)
				break
			}
			data = rest

			count++
			r.pipeline <- messageErrorPair{msg, parseErr}
			if len(rest) == 0 {
				break
			}
		}

		done = len(data)
		buffer = buffer[lenRead-done:]

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

// Receive an incoming syslog message and a possible error that occured during the decoding of this syslog message
func (r *Receiver) Receive() (MessageImmutable, error, bool) {
	pair, end := <-r.pipeline

	return pair.m, pair.e, end
}

// End terminate the log_collector goroutine
func (r *Receiver) End() {
	close(r.end)
}
