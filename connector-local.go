package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"net"
)

type (
	localConn struct {
		network, address string
	}

	unixgram struct {
		addr *Addr
		c    *net.UnixConn
	}
)

// LocalConnector is a dialer that forward to a local RFC5424 syslog receiver
func LocalConnector(network, address string) Connector {
	return &localConn{network, address}
}

// Connect is part of implementation of (Connector interface)[#Connector]
func (c *localConn) Connect() (WriteCloser, error) {
	if c.address != "" && c.network != "" {
		return c.localWriteCloser(net.DialUnix(c.network, nil, &net.UnixAddr{Name: c.address, Net: c.network}))
	}

	return c.localWriteCloser(c.osGuessConnnector())
}

func (c *localConn) localWriteCloser(conn *net.UnixConn, err error) (WriteCloser, error) {
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(0)
	conn.CloseRead()

	if c.network == "unixgram" {
		//		return conn, nil
		return unixgram{&Addr{c.network, c.address}, conn}, nil
	}

	return conn, nil
}

func (c unixgram) Close() error {
	return c.c.Close()
}

func (c unixgram) Write(d []byte) (n int, err error) {
	return c.c.Write(d)
}
