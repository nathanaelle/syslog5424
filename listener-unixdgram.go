package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"net"
	"os"
)

type (
	unixgram_receiver struct {
		listener *net.UnixConn
		accepted bool
		end      chan struct{}
	}

	fake_conn struct {
		end   chan struct{}
		rbuff [1 << 16]byte
		buff  []byte
		conn  *net.UnixConn
	}
)

func UnixgramListener(address string) (Listener, error) {
	var err error

	r := new(unixgram_receiver)
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

	r.listener.SetWriteBuffer(0)
	r.listener.CloseWrite()
	r.listener.SetReadBuffer(readBuffer)

	return r, nil
}

func (r *unixgram_receiver) Close() error {
	close(r.end)
	return r.listener.Close()
}

// mimic an Accept
func (r *unixgram_receiver) Accept() (DataReader, error) {
	if r.accepted {
		<-r.end
		return nil, errors.New("end")
	}

	r.accepted = true

	fc := &fake_conn{
		end:  r.end,
		conn: r.listener,
	}

	return fc, nil
}

func (r *fake_conn) RemoteAddr() net.Addr {
	return r.conn.RemoteAddr()
}

func (r *fake_conn) Close() error {
	return nil
}

func (r *fake_conn) Read(data []byte) (int, error) {
	if len(r.buff) == 0 {
		s, _, err := r.conn.ReadFrom(r.rbuff[:])
		if err != nil {
			return s, err
		}

		r.buff = make([]byte, s)
		copy(r.buff, r.rbuff[0:s])
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
