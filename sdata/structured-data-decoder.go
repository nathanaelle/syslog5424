package sdata // import "github.com/nathanaelle/syslog5424/sdata"

import (
	"reflect"
	"runtime"
)

type (
	decodeState struct {
	}

	InvalidUnmarshalError struct {
		Type reflect.Type
	}
)

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "sd5424: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "sd5424: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "sd5424: Unmarshal(nil " + e.Type.String() + ")"
}

func UnmarshalSD(data []byte, sd_i interface{}) error {
	var d decodeState

	sd := reflect.ValueOf(sd_i)
	if sd.Kind() != reflect.Ptr || sd.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(sd_i)}
	}

	return d.unmarshal(sd_i)
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	return nil
}
