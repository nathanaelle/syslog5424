package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
)

type (
	tcp_receiver struct {
		network   string
		address   string
		listener  *net.TCPListener
		transport Transport
		pipeline  chan []byte
	}
)

func tcp_coll(network, address string) (Listener, error) {
	var err error

	r := new(tcp_receiver)
	r.network = network
	r.address = address

	laddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	r.listener, err = net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *tcp_receiver) Accept() (net.Conn, error) {
	conn, err := r.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetWriteBuffer(0)
	conn.SetReadBuffer(readBuffer)

	return conn, nil
}

func (r *tcp_receiver) Close() error {
	return r.listener.Close()
}

func (r *tcp_receiver) Addr() net.Addr {
	return r.listener.Addr()
}
