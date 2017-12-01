package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"os"
)

// StdioConnector returns a Connector that only forward to stderr: or stdout:
func StdioConnector(addr string) Connector {
	if addr == "" {
		return InvalidConnector{ErrorEmptyNetworkAddress}
	}

	switch addr {
	case "stderr:":
		return ConnectorFunc(func() (WriteCloser, error) {
			return os.Stderr, nil
		})

	case "stdout:":
		return ConnectorFunc(func() (WriteCloser, error) {
			return os.Stdout, nil
		})
	}

	// TODO implement file logging here

	return InvalidConnector{ErrorInvalidAddress}
}
