package sdata // import "github.com/nathanaelle/syslog5424/v2/sdata"

import (
	"fmt"
	"reflect"
)

type (
	// MarshalerError …
	MarshalerError struct {
		Type reflect.Type
		Err  error
	}

	// UnsupportedTypeError …
	UnsupportedTypeError struct {
		Type reflect.Type
	}

	// InvalidValueError …
	InvalidValueError struct {
		Value reflect.Value
	}

	// InvalidUnmarshalError …
	InvalidUnmarshalError struct {
		Type reflect.Type
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

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "sd5424: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "sd5424: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "sd5424: Unmarshal(nil " + e.Type.String() + ")"
}
