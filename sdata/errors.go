package sdata // import "github.com/nathanaelle/syslog5424/sdata"

import (
	"fmt"
	"reflect"
)

type (
	MarshalerError struct {
		Type reflect.Type
		Err  error
	}

	UnsupportedTypeError struct {
		Type reflect.Type
	}

	InvalidValueError struct {
		Value reflect.Value
	}
)

func (e *UnsupportedTypeError) Error() string {
	return "StructuredData: unsupported type: " + e.Type.String()
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("StructuredData: unsupported type: %#v", e.Value)
}

func (e *MarshalerError) Error() string {
	return "StructuredData: error calling Marshal5424 type " + e.Type.String() + ": " + e.Err.Error()
}
