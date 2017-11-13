package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
)

type (
	tcp_receiver struct {
		listener  *net.TCPListener
	}
)

func TCPListener(network, address string) (Listener, error) {
	var err error

	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, ErrorInvalidNetwork
	}


	laddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	r := new(tcp_receiver)
	r.listener, err = net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *tcp_receiver) Accept() (DataReader, error) {
	conn, err := r.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetLinger(-1)
	conn.CloseWrite()
	conn.SetReadBuffer(readBuffer)

	return conn, nil
}

func (r *tcp_receiver) Close() error {
	return r.listener.Close()
}
