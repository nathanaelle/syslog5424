package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
	"os"
)

type (
	unix_receiver struct {
		listener  *net.UnixListener
	}
)

// careful : a previous unused socket may be removed and recreated
func UnixListener(address string) (Listener, error) {
	var err error

	r := new(unix_receiver)
	r.listener, err = net.ListenUnix("unix", &net.UnixAddr{address, "unix"})

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

		r.listener, err = net.ListenUnix("unix", &net.UnixAddr{address, "unix"})
	}

	return r, nil
}

func (r *unix_receiver) Accept() (DataReader, error) {
	conn, err := r.listener.AcceptUnix()
	if err != nil {
		return nil, err
	}
	conn.SetWriteBuffer(0)
	conn.CloseWrite()
	conn.SetReadBuffer(readBuffer)

	return conn, nil
}

func (r *unix_receiver) Close() error {
	return r.listener.Close()
}
