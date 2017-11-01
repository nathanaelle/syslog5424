package syslog5424 // import "github.com/nathanaelle/syslog5424"

import ()

type (
	local_conn struct {
		fd_conn
	}
)

// dialer that forward to a local RFC5424 syslog receiver
func local_dial(network, address string) Conn {
	s := new(local_conn)

	s.address = address
	s.network = network
	s.writer = new_buffer(1<<10, buffer_write, nil)

	return s
}

func (c *local_conn) Close() error {
	if c.writer != nil {
		c.writer.Flush()
		return c.writer.Close()
	}
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

func (c *local_conn) Redial() error {
	conn, err := c.os_redial()
	if err != nil {
		return err
	}

	c.writer.SetConn(conn)
	return nil
}
