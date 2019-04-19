package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"net"
)

type (
	tcpReceiver struct {
		listener *net.TCPListener
	}
)

// TCPListener create a TCP Listener
func TCPListener(network, address string) (Listener, error) {
	var err error

	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, ErrInvalidNetwork
	}

	laddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	r := new(tcpReceiver)
	r.listener, err = net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *tcpReceiver) Accept() (DataReader, error) {
	conn, err := r.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetLinger(-1)
	conn.CloseWrite()
	conn.SetReadBuffer(readBuffer)

	return conn, nil
}

func (r *tcpReceiver) Close() error {
	return r.listener.Close()
}
