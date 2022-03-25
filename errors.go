package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"errors"
	"fmt"
)

type (
	// ParseError describe an error in parsing
	ParseError struct {
		Buffer  []byte
		Pos     int
		Message string
	}
)

var (
	ErrBufferClose  = errors.New("error in syslog5424 at buffer.Close()")
	ErrNoConnection = errors.New("No Connection established")

	ErrPos0                = errors.New("Pos 0 Found")
	ErrPosNotFound         = errors.New("Pos Not Found")
	ErrImpossible          = errors.New("NO ONE EXPECT THE RETURN OF SPANISH INQUISITION")
	ErrInvalidNetwork      = errors.New("Invalid Network")
	ErrInvalidAddress      = errors.New("Invalid Address")
	ErrEmptyNetworkAddress = errors.New("Empty Network or Address")

	ErrTransportIncomplete = errors.New("Transport : Incomplete Message")
	ErrTransportNoHeader   = errors.New("T_RFC5425 Split: no header len")
	ErrTransportInvHeader  = errors.New("T_RFC5425 Split: invalid header len")
)

func (pe ParseError) Error() string {
	return fmt.Sprintf("{%q} %d %s", pe.Buffer, pe.Pos, pe.Message)
}
