// +build linux freebsd netbsd openbsd dragonfly

package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
)

func (c *local_conn) osGuessConnnector() (*net.UnixConn, error) {
	logTypes := []string{"unix", "unixgram"}
	logPaths := []string{"/var/run/syslog", "/var/run/log", "/dev/log"}

	if c.address != "" {
		for _, network := range logTypes {
			conn, err := net.Dial(network, c.address)
			if err == nil {
				c.network = network
				return conn, nil
			}
		}
		return nil, ERR_NOCONN
	}

	for _, network := range logTypes {
		for _, path := range logPaths {
			conn, err := net.Dial(network, path)
			if err == nil {
				c.network = network
				c.address = path
				return conn, nil
			}
		}
	}
	return nil, ERR_NOCONN
}
