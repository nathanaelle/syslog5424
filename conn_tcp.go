package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
)


type (
	tcp_conn struct {
		fd_conn
	}
)


// dialer that forward to a local RFC5424 syslog receiver
func tcp_dial(network, address string) Conn {
	s := new(tcp_conn)

	s.address	= address
	s.network	= network
	s.writer	= new_buffer(1<<12, buffer_write, nil)

	return s
}



func (c *tcp_conn) Close() error {
	if c.writer != nil {
		return c.writer.Close()
	}
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}


func (c *tcp_conn) Redial() (error) {
	conn, err := net.Dial(c.network, c.address)
	if err != nil {
		return err
	}

	c.writer.SetConn( conn )
	return nil
}
