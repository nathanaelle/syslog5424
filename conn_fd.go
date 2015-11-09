package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"os"
)


type (
	fd_conn struct {
		address 	string
		network		string
		reader		*buffer
		writer		*buffer
	}
)


// dialer that only forward to stderr
func stdio_dial(addr string) Conn {
	s := new(fd_conn)

	s.address	= addr
	s.network	= "stdio"

	switch addr {
	case "stderr":
		s.writer	= new_buffer(1<<12, buffer_write, os.Stderr)

	case "stdout":
		s.writer	= new_buffer(1<<12, buffer_write, os.Stdout)

	// TODO implement file logging here
	default:
		return nil
	}

	return	s
}


func (c *fd_conn) Write(data []byte) (n int, err error) {
	var t_n int
	for n < len(data) {
		t_n, err = c.writer.Write(data[n:])
		n += t_n
		if err != nil {
			return
		}
	}
	return
}


func (c *fd_conn) Read(data []byte) (int, error) {
	return c.reader.Read(data)
}


func (c *fd_conn) Redial() error {
	switch c.address {
	case "stderr","stdout":
		return nil

	// TODO implement file logging here
	default:
		return nil
	}
}


func (c *fd_conn) Flush() error {
	if c.writer != nil {
		return c.writer.Flush()
	}
	if c.reader != nil {
		return c.reader.Flush()
	}
	return nil
}


func (c *fd_conn) Close() error {
	c.Flush()

	switch c.address {
	case "stderr","stdout":
		return nil

	// TODO implement file logging here
	default:
		return nil
	}
}
