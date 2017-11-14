package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"time"
)

type (
	Dialer struct {
		// delay to flush the queue
		FlushDelay time.Duration
	}
)

// Dial opens a connection to the syslog daemon
// network can be "stdio", "unix", "unixgram", "tcp", "tcp4", "tcp6"
// used Transport is the "common" transport for the network.
// FlushDelay is preset to 500ms
func Dial(network, address string) (*Sender, <-chan error, error) {
	return (Dialer{
		FlushDelay: 500 * time.Millisecond,
	}).Dial(network, address, nil)
}

// Dial opens a connection to the syslog daemon
// network can be "stdio", "unix", "unixgram", "tcp", "tcp4", "tcp6"
// Transport can be nil.
// if Transport is nil the "common" transport for the wished network is used.
//
// the returned `<-chan error` is used to collect errors than may occur in goroutine
func (d Dialer) Dial(network, address string, t Transport) (*Sender, <-chan error, error) {
	var ticker <-chan time.Time
	var c Connector

	switch {
	case d.FlushDelay <= time.Millisecond:
		// less than 1ms => disable auto flush
		ticker = make(chan time.Time)

	default:
		ticker = time.Tick(d.FlushDelay)
	}

	switch network {
	case "stdio":
		if t == nil {
			t = T_LFENDED
		}
		c = StdioConnector(address)

	case "local":
		if t == nil {
			t = T_ZEROENDED
		}
		c = LocalConnector("", address)

	case "unix", "unixgram":
		if t == nil {
			t = T_ZEROENDED
		}
		c = LocalConnector(network, address)

	case "tcp", "tcp6", "tcp4":
		if t == nil {
			t = T_LFENDED
		}
		c = TCPConnector(network, address)

	default:
		return nil, nil, errors.New("unknown network for Dial : " + network)
	}

	if c == nil {
		return nil, nil, ErrorNoConnecion
	}

	sndr, chan_err := NewSender(c, t, ticker)
	return sndr, chan_err, nil
}
