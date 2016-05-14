// xx+build linux freebsd netbsd openbsd dragonfly

package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"net"
	"errors"
)

func (c *local_conn) os_redial() (io.ReadWriteCloser, error) {
	logTypes := []string{"unix","unixgram"}
	logPaths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}

	if c.address != "" && c.network != "" {
		return net.Dial(c.network, c.address)
	}

	if c.address != "" {
		for _, network := range logTypes {
			conn, err := net.Dial(network, c.address)
			if err == nil {
				c.network = network
				return conn,nil
			}
		}
		return nil,errors.New("no connection established")
	}

	for _, network := range logTypes {
		for _, path := range logPaths {
			conn, err := net.Dial(network, path)
			if err == nil {
				c.network = network
				c.address = path
				return conn,nil
			}
		}
	}
	return nil,errors.New("no connection established")
}
