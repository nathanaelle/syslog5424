package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"crypto/tls"
	"io"
	"net"
	"os"
	"strconv"
)

type (
	Transport int

	Conn interface {
		io.WriteCloser
		Redial() error
		Send(Message)
		End()
	}

	local_conn struct {
		fd_conn
		address string
		network string
	}

	tcp_conn struct {
		local_conn
	}

	tls_conn struct {
		tcp_conn
	}

	fd_conn struct {
		pipeline chan Message
		conn     io.WriteCloser
	}
)

const (
	T_ZEROENDED Transport = iota
	T_LFENDED
	T_RFC5426
)

func (t Transport) Encoder() func([]byte) []byte {
	switch t {
	case T_ZEROENDED:
		return func(d []byte) []byte {
			return append(d, 0)
		}

	case T_LFENDED:
		return func(d []byte) []byte {
			return append(d, '\n')
		}

	case T_RFC5426:
		return func(d []byte) []byte {
			l := len(d)
			h := []byte(strconv.Itoa(l))
			ret := make([]byte, l+len(h)+1)
			copy(ret[0:len(h)], h[:])
			ret[len(h)] = ' '
			copy(ret[len(h)+1:], d[:])
			return ret
		}
	}

	return func([]byte) []byte {
		panic("unknown Transport Encoder")
	}
}

func task_logger(pipeline <-chan Message, output Conn, encode func([]byte) []byte) {
	output.Redial()

	for {
		select {
		case msg, opened := <-pipeline:
			if !opened {
				break
			}

			w := 0
			d := encode(msg.Marshal5424())
			for w < len(d) {
				t_w, err := output.Write(d[w:])
				w += t_w
				if err != nil {
					output.Redial()
					w = 0
				}
			}
		}
	}
	output.Close()
}

// Dial opens a connection to the syslog daemon
// network can be "local", "unixgram", "tcp", "tcp4", "tcp6"
func Dial(network, address string, t Transport, queue_len int) Conn {
	var c Conn

	if queue_len < 100 {
		queue_len = 100
	}
	pipeline := make(chan Message, queue_len)

	switch network {
	case "stderr":
		c = stderr_dial("", pipeline)

	case "local", "unixgram":
		c = local_dial(address, pipeline)

	case "tcp", "tcp6", "tcp4":
		c = tcp_dial(network, address, pipeline)
	default:
		return nil
	}

	go task_logger(pipeline, c, t.Encoder())

	return c
}

// TLSDial opens a connection to the syslog daemon
// network can be "local", "unixgram", "tcp", "tcp4", "tcp6"
// TODO write the code
func TLSDial(network, address string, t Transport, o tls.Config) Conn {
	var c Conn
	pipeline := make(chan Message, 100)

	switch network {
	case "local", "unixgram":
		c = local_dial(address, pipeline)

	case "tcp", "tcp6", "tcp4":
		c = tcp_dial(network, address, pipeline)
	default:
		return nil
	}

	go task_logger(pipeline, c, t.Encoder())

	return c
}

// dialer that only forward to stderr
func stderr_dial(_ string, pipeline chan Message) Conn {
	return &fd_conn{
		pipeline: pipeline,
		conn:     os.Stderr,
	}
}

// dialer that forward to a local RFC5424 syslog receiver
func local_dial(address string, pipeline chan Message) Conn {
	c := new(local_conn)
	c.address = address
	c.pipeline = pipeline
	return c
}

// dialer that forward to a local RFC5424 syslog receiver
func tcp_dial(network, address string, pipeline chan Message) Conn {
	c := new(tcp_conn)
	c.address = address
	c.network = network
	c.pipeline = pipeline
	return c
}

func (c *fd_conn) Write(data []byte) (n int, err error) {
	var t_n int
	for n < len(data) {
		t_n, err = c.conn.Write(data[n:])
		n += t_n
		if err != nil {
			return
		}
	}
	return
}

func (c *fd_conn) Redial() error {
	return nil
}

func (c *fd_conn) Send(m Message) {
	c.pipeline <- m
}

func (c *fd_conn) End() {
	close(c.pipeline)
}

func (c *fd_conn) Close() error {
	return nil
}

func (c *local_conn) Close() error {
	return c.conn.Close()
}

func (c *tcp_conn) Redial() (err error) {
	c.conn, err = net.Dial(c.network, c.address)
	return
}
