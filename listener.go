package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"net"
	"bufio"
	"errors"
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
	}


	Receiver struct {
		listener	Listener
		transport	Transport
		pipeline	chan []byte
		end		chan struct{}
	}
)


func Collect(network,address string) (*Receiver,error) {
	return (Collector{
		QueueLen:	100,
	}).Collect(network, address, nil)
}


func (d Collector) Collect(network,address string, t Transport) (*Receiver,error) {
	var pipeline	chan []byte
	var c		Listener
	var err		error

	switch network {
	case "unix":
		if t == nil {
			t = new(T_ZEROENDED)
		}
		c,err = unix_coll(network, address)

	case "unixgram":
		if t == nil {
			t = new(T_ZEROENDED)
		}
		c,err = unixgram_coll(network, address)

	case "tcp", "tcp6", "tcp4":
		if t == nil {
			t = new(T_LFENDED)
		}
		c,err = tcp_coll(network, address)

	default:
		return nil, errors.New("unknown network for Collector : "+network)
	}

	if err != nil {
		return nil, err
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

	return NewReceiver(c, pipeline, t),nil
}


func NewReceiver(listener Listener, pipeline chan []byte, t Transport) (*Receiver) {
	r	:= &Receiver {
		listener:	listener,
		pipeline:	pipeline,
		transport:	t,
		end:		make(chan struct{}),
	}

	go r.run_queue()

	return	r
}


func (r *Receiver) run_queue() {
	defer	r.listener.Close()
	defer	close(r.pipeline)

	for {
		select {
		case <-r.end:
			return

		default:
			conn, err := r.listener.Accept()
			if err != nil {
				panic(err)
			}

			go r.tokenize(new_buffer(1<<18, buffer_read, conn))
		}
	}

}

func (r *Receiver) tokenize(conn io.ReadWriteCloser) {
	scan	:= bufio.NewScanner(conn)
	scan.Split(r.transport.Split)

	for scan.Scan() {
		r.pipeline <- scan.Bytes()
	}

	conn.Close()
}


func (r *Receiver) ReceiveRaw() ([]byte, bool) {
	b,end	:= <- r.pipeline
	return b, end
}



func (r *Receiver) Receive() (Message, error, bool) {
	b,end	:= r.ReceiveRaw()
	msg,err	:= Parse(b)

	return msg, err, end
}


// terminate the log_collector goroutine
func (r *Receiver) End() {
	close(r.end)
}
