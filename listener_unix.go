package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
	"os"
)

type (
	unix_receiver struct {
		network   string
		address   string
		listener  net.Listener
		transport Transport
		pipeline  chan []byte
	}
)

func unix_coll(_, address string) (Listener, error) {
	var err error

	r := new(unix_receiver)
	r.network = "unix"
	r.address = address

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

func (r *unix_receiver) Accept() (net.Conn, error) {
	return r.listener.Accept()
}

func (r *unix_receiver) Close() error {
	return r.listener.Close()
}

func (r *unix_receiver) Addr() net.Addr {
	return &Addr{r.network, r.address}
}
