package sdata // import "github.com/nathanaelle/syslog5424/v2/sdata"

import (
	"reflect"
	"runtime"
)

type (
	decodeState struct {
	}
)

// UnmarshalSD decode a []byte to a golang struct
func UnmarshalSD(data []byte, v interface{}) error {
	var d decodeState

	sd := reflect.ValueOf(v)
	if sd.Kind() != reflect.Ptr || sd.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return d.unmarshal(v)
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
