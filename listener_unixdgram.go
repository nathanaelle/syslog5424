package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"net"
	"os"
	"time"
)

type (
	unixgram_receiver struct {
		network  string
		address  string
		listener *net.UnixConn
		accepted bool
		end      chan struct{}
	}

	fake_conn struct {
		addr  *Addr
		buff  []byte
		queue chan []byte
		end   chan struct{}
	}
)

func unixgram_coll(_, address string) (Listener, error) {
	var err error

	r := new(unixgram_receiver)
	r.network = "unixgram"
	r.address = address
	r.end = make(chan struct{})

	r.listener, err = net.ListenUnixgram("unixgram", &net.UnixAddr{address, "unixgram"})
	for err != nil {
		switch err.(type) {
		case *net.OpError:
			if err.(*net.OpError).Err.Error() != "bind: address already in use" {
				return nil, err
			}

		default:
			return nil, err
		}

		if _, r_err := os.Stat(address); r_err != nil {
			return nil, err
		}
		os.Remove(address)

		r.listener, err = net.ListenUnixgram("unixgram", &net.UnixAddr{address, "unixgram"})
	}
	return r, nil
}

func (r *unixgram_receiver) Close() error {
	close(r.end)
	return r.listener.Close()
}

func (r *unixgram_receiver) Addr() net.Addr {
	return &Addr{r.network, r.address}
}

// mimic an Accept
func (r *unixgram_receiver) Accept() (net.Conn, error) {
	if r.accepted {
		<-r.end
		return nil, errors.New("end")
	}

	r.accepted = true

	fc := &fake_conn{
		addr:  &Addr{r.network, r.address},
		queue: make(chan []byte, 1000),
		end:   r.end,
	}

	go fc.run_queue(r.listener)

	return fc, nil
}

func (r *fake_conn) LocalAddr() net.Addr {
	return r.addr
}

func (r *fake_conn) RemoteAddr() net.Addr {
	return r.addr
}

func (r *fake_conn) SetDeadline(_ time.Time) error {
	return nil
}

func (r *fake_conn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (r *fake_conn) SetWriteDeadline(_ time.Time) error {
	return nil
}

func (c *fake_conn) Redial() error {
	return nil
}

func (c *fake_conn) Flush() error {
	return nil
}

func (r *fake_conn) Close() error {
	return nil
}

func (r *fake_conn) Write(data []byte) (int, error) {
	return len(data), nil
}

func (r *fake_conn) Read(data []byte) (int, error) {
	if len(r.buff) == 0 {
		r.buff = <-r.queue
	}

	l_r := len(r.buff)
	l_d := len(data)
	if l_d <= l_r {
		copy(data[:], r.buff[0:l_d])
		r.buff = r.buff[l_d:]
		return l_d, nil
	}

	copy(data[0:l_r], r.buff[:])
	r.buff = nil

	return l_r, nil
}

func (r *fake_conn) run_queue(conn *net.UnixConn) {
	defer conn.Close()

	for {
		select {
		case <-r.end:
			return

		default:
			buffer := make([]byte, 1<<16)

			conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
			s, _, err := conn.ReadFrom(buffer)
			switch t_err := err.(type) {
			case nil:
			case net.Error:
				if !t_err.Timeout() {
					panic(err)
				}
			default:
				panic(err)
			}

			if s > 0 {
				r.queue <- buffer[0:s]
			}

		}
	}

}
