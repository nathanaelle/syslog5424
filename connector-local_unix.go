// +build darwin linux freebsd netbsd openbsd dragonfly

package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
)

func (c *localConn) osGuessConnnector() (*net.UnixConn, error) {
	logTypes := []string{"unix", "unixgram"}
	logPaths := []string{"/var/run/syslog", "/var/run/log", "/dev/log"}

	if c.address != "" {
		for _, network := range logTypes {
			conn, err := net.DialUnix(network, nil, &net.UnixAddr{c.address, network})
			if err == nil {
				c.network = network
				return conn, nil
			}
		}
		return nil, ErrorNoConnecion
	}

	for _, network := range logTypes {
		for _, path := range logPaths {
			conn, err := net.DialUnix(network, nil, &net.UnixAddr{path, network})
			if err == nil {
				c.network = network
				c.address = path
				return conn, nil
			}
		}
	}
	return nil, ErrorNoConnecion
}
