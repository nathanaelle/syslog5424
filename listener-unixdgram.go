package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"net"
	"os"
)

type (
	unixgramReceiver struct {
		listener *net.UnixConn
		accepted bool
		end      chan struct{}
	}

	fakeConn struct {
		end   chan struct{}
		rbuff [1 << 16]byte
		buff  []byte
		conn  *net.UnixConn
	}
)

func UnixgramListener(address string) (Listener, error) {
	var err error

	r := new(unixgramReceiver)
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

		if _, osErr := os.Stat(address); osErr != nil {
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

func (r *unixgramReceiver) Close() error {
	close(r.end)
	return r.listener.Close()
}

// mimic an Accept
func (r *unixgramReceiver) Accept() (DataReader, error) {
	if r.accepted {
		<-r.end
		return nil, errors.New("end")
	}

	r.accepted = true

	fc := &fakeConn{
		end:  r.end,
		conn: r.listener,
	}

	return fc, nil
}

func (r *fakeConn) RemoteAddr() net.Addr {
	return r.conn.RemoteAddr()
}

func (r *fakeConn) Close() error {
	return nil
}

func (r *fakeConn) Read(data []byte) (int, error) {
	if len(r.buff) == 0 {
		s, _, err := r.conn.ReadFrom(r.rbuff[:])
		if err != nil {
			return s, err
		}

		r.buff = make([]byte, s)
		copy(r.buff, r.rbuff[0:s])
	}

	lenBuff := len(r.buff)
	lenData := len(data)
	if lenData <= lenBuff {
		copy(data[:], r.buff[0:lenData])
		r.buff = r.buff[lenData:]
		return lenData, nil
	}

	copy(data[0:lenBuff], r.buff[:])
	r.buff = nil

	return lenBuff, nil
}
