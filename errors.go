package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"fmt"
)

type (
	ParseError struct {
		Buffer  []byte
		Pos     int
		Message string
	}
)

var (
	ErrorBufferClose error = errors.New("error in syslog5424 at buffer.Close()")
	ErrorNoConnecion error = errors.New("No Connection established")

	ErrorPos0                error = errors.New("Pos 0 Found")
	ErrorPosNotFound         error = errors.New("Pos Not Found")
	ErrorImpossible          error = errors.New("NO ONE EXPECT THE RETURN OF SPANISH INQUISITION")
	ErrorInvalidNetwork      error = errors.New("Invalid Network")
	ErrorInvalidAddress      error = errors.New("Invalid Address")
	ErrorEmptyNetworkAddress error = errors.New("Empty Network or Address")

	ERR_TRANSPORT_INCOMPLETE error = errors.New("Transport : Incomplete Message")
	ERR_TRANSPORT_NOHEADER   error = errors.New("T_RFC5425 Split: no header len")
	ERR_TRANSPORT_INVHEADER  error = errors.New("T_RFC5425 Split: invalid header len")
)

func (pe ParseError) Error() string {
	return fmt.Sprintf("{%q} %d %s", pe.Buffer, pe.Pos, pe.Message)
}
