package syslog5424 // import "github.com/nathanaelle/syslog5424"

type (
	InvalidConnector struct {
		err error
	}
)

func (ic InvalidConnector) Connect() (conn WriteCloser, err error) {
	return nil, ic.err
}
