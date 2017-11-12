// +build !linux,!freebsd,!netbsd,!openbsd,!dragonfly

package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"net"
)

func (c *local_conn) osGuessConnnector() (*net.UnixConn, error) {
	return nil, ErrorNoConnecion
}
