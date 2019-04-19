// +build !darwin,!linux,!freebsd,!netbsd,!openbsd,!dragonfly

package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"net"
)

func (c *localConn) osGuessConnnector() (*net.UnixConn, error) {
	return nil, ErrorNoConnecion
}
