package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"net"
	"os"
)

type (
	unixReceiver struct {
		listener *net.UnixListener
	}
)

// UnixListener create a UNIX Listener
// careful : a previous unused socket may be removed and recreated
func UnixListener(address string) (Listener, error) {
	var err error

	r := new(unixReceiver)
	uAddr, err := net.ResolveUnixAddr("unix", address)
	if err != nil {
		return nil, err
	}

	r.listener, err = net.ListenUnix("unix", uAddr)

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

		r.listener, err = net.ListenUnix("unix", uAddr)
	}

	return r, nil
}

func (r *unixReceiver) Accept() (DataReader, error) {
	conn, err := r.listener.AcceptUnix()
	if err != nil {
		return nil, err
	}
	conn.SetWriteBuffer(0)
	conn.CloseWrite()
	conn.SetReadBuffer(readBuffer)

	return conn, nil
}

func (r *unixReceiver) Close() error {
	return r.listener.Close()
}
