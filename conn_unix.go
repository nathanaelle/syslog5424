package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"net"
)

func (c *local_conn) Redial() (err error) {
	logTypes := []string{"unixgram", "unix"}
	logPaths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}

	if c.address != "" && c.network != "" {
		c.conn, err = net.Dial(c.network, c.address)
		return
	}

	if c.address != "" {
		for _, network := range logTypes {
			conn, err := net.Dial(network, c.address)
			if err == nil {
				c.conn = conn
				c.network = network
				return nil
			}
		}
		return errors.New("no connection established")
	}

	for _, network := range logTypes {
		for _, path := range logPaths {
			conn, err := net.Dial(network, path)
			if err == nil {
				c.conn = conn
				c.network = network
				c.address = path
				return nil
			}
		}
	}
	return errors.New("no connection established")
}
