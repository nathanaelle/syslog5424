// +build darwin linux freebsd netbsd openbsd dragonfly

package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"net"
)

func (c *localConn) osGuessConnnector() (*net.UnixConn, error) {
	logTypes := []string{"unix", "unixgram"}
	logPaths := []string{"/var/run/syslog", "/var/run/log", "/dev/log"}

	if c.address != "" {
		for _, network := range logTypes {
			uAddr, err := net.ResolveUnixAddr(network, c.address)
			if err != nil {
				continue
			}
			conn, err := net.DialUnix(network, nil, uAddr)
			if err == nil {
				c.network = network
				return conn, nil
			}
		}
		return nil, ErrNoConnecion
	}

	for _, network := range logTypes {
		for _, path := range logPaths {
			uAddr, err := net.ResolveUnixAddr(network, path)
			if err != nil {
				continue
			}
			conn, err := net.DialUnix(network, nil, uAddr)
			if err == nil {
				c.network = network
				c.address = path
				return conn, nil
			}
		}
	}
	return nil, ErrNoConnecion
}
