package syslog5424 // import "github.com/nathanaelle/syslog5424"

type (
	// InvalidConnector is a connector that always return an error
	InvalidConnector struct {
		err error
	}
)

// Connect is part of implementation of (Connector interface)[#Connector]
func (ic InvalidConnector) Connect() (conn WriteCloser, err error) {
	return nil, ic.err
}
