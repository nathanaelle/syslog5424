package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
)

type (
	encodeState struct {
		bytes.Buffer
	}

	encoderFunc  func(*encodeState, reflect.Value, []byte)
	stringValues []reflect.Value
)

const (
	caseMask = ^byte(0x20) // Mask to ignore case in ASCII.
)

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

func MarshalSD(sd interface{}) []byte {
	e := &encodeState{}
	err := e.marshal(sd)
	if err != nil {
		return []byte{}
	}
	return e.Bytes()
}

func (e *encodeState) marshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case runtime.Error, string:
				panic(r)
			}
			err = r.(error)
		}
	}()

	e.reflectValue(reflect.ValueOf(v))
	return nil
}

func (e *encodeState) PrefixedWrite(prefix, data []byte) {
	switch len(prefix) == 0 {
	case true:
		e.Write(data)

	case false:
		e.WriteByte(' ')
		e.Write(prefix)
		e.WriteByte('=')
		e.WriteByte('"')
		for _, d := range data {
			if d == '"' || d == '\\' || d == ']' {
				e.WriteByte('\\')
			}
			e.WriteByte(d)
		}
		e.WriteByte('"')
	}
}

func (e *encodeState) HeaderWrite(typ string, v reflect.Value) {
	e.WriteByte('[')
	e.Write([]byte(typ))

	if v.Type().Implements(sdpenType) {
		e.WriteByte('@')
		pen := strconv.FormatUint(v.Interface().(SDPEN).GetPEN(), 10)
		e.Write([]byte(pen))
	}
}

func (e *encodeState) error(err error) {
	panic(err)
}

func (e *encodeState) reflectValue(v reflect.Value) {
	newTypeEncoder(v.Type(), true)(e, v, []byte{})
}

func unsupportedTypeEncoder(e *encodeState, v reflect.Value, _ []byte) {
	e.error(&UnsupportedTypeError{v.Type()})
}

func invalidValueEncoder(e *encodeState, v reflect.Value, _ []byte) {
	e.error(&InvalidValueError{v})
}

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return newTypeEncoder(v.Type(), false)
}

func newTypeEncoder(t reflect.Type, root bool) encoderFunc {
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}

	switch root {
	case true:
		switch t.Kind() {
		case reflect.Struct:
			return newStructEncoder(t)

		case reflect.Map:
			return newMapEncoder(t)

		default:
			if t.Kind() != reflect.Ptr {
				return newTypeEncoder(reflect.PtrTo(t), true)
			}
		}

	case false:
		switch t.Kind() {
		case reflect.Ptr:
			return newMaybeEncoder(t)

		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			return encodeur_flemard

		case reflect.Interface:
			return newTypeEncoder(t.Elem(), false)

		case reflect.Slice, reflect.Array:
			return newListEncoder(t)

		}
	}
	return unsupportedTypeEncoder
}

// encoder that delegate to the SDMarshaler the encoding
func marshalerEncoder(e *encodeState, v reflect.Value, prefix []byte) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return
	}

	m := v.Interface().(SDMarshaler)
	b, err := m.Marshal5424()
	if err != nil {
		e.error(&MarshalerError{v.Type(), err})
		return
	}
	e.PrefixedWrite(prefix, b)
}

// trivial encoder for string, int, float, bool, ... based on Sprintf
func encodeur_flemard(e *encodeState, v reflect.Value, prefix []byte) {
	e.PrefixedWrite(prefix, []byte(fmt.Sprintf("%v", v.Interface())))
}

// encoder for pointer on type
// this encoder act like a "?" pattern in a regexp
func newMaybeEncoder(t reflect.Type) encoderFunc {
	return func(e *encodeState, v reflect.Value, prefix []byte) {
		if v.IsNil() {
			return
		}
		newTypeEncoder(t.Elem(), false)(e, v.Elem(), prefix)
	}
}

// encoder for list, array, slice of a type
// this encoder acts like a "*" pattern in a regexp
func newListEncoder(t reflect.Type) encoderFunc {
	return func(e *encodeState, v reflect.Value, prefix []byte) {
		if v.IsNil() {
			return
		}
		enc := newTypeEncoder(t.Elem(), false)
		n := v.Len()
		for i := 0; i < n; i++ {
			enc(e, v.Index(i), prefix)
		}
	}
}

// encoder for map
func newMapEncoder(t reflect.Type) encoderFunc {
	if t.Key().Kind() != reflect.String {
		return unsupportedTypeEncoder
	}
	enc := newTypeEncoder(t.Elem(), false)

	return func(e *encodeState, v reflect.Value, _ []byte) {
		if v.IsNil() {
			return
		}

		e.HeaderWrite(t.Name(), v)

		var sv stringValues = v.MapKeys()
		sort.Sort(sv)
		for _, k := range sv {
			enc(e, v.MapIndex(k), []byte(k.String())[:])
		}
		e.WriteByte(']')
	}
}

func newStructEncoder(t reflect.Type) encoderFunc {
	n := t.NumField()
	fields := make([]reflect.StructField, n)
	for i := 0; i < n; i++ {
		fields[i] = t.Field(i)
	}
	return func(e *encodeState, v reflect.Value, _ []byte) {
		e.HeaderWrite(t.Name(), v)

		for i := 0; i < n; i++ {
			sf := fields[i]
			f := v.Field(i)
			name, mandatory := getName(sf)
			if !f.IsValid() || isEmptyValue(f) || len(name) == 0 {
				if mandatory {
					// TODO this must be an error
					panic("missing mandatory field")
				}
				continue
			}
			newTypeEncoder(sf.Type, false)(e, f, name)
		}
		e.WriteByte(']')
	}
}

func getName(field reflect.StructField) ([]byte, bool) {
	name := field.Name
	tag := field.Tag.Get("sd5424")

	if tag == "-" {
		return []byte{}, false
	}
	t_name, opts := parseTag(tag)
	mandatory := opts.Contains("mandatory")
	if isValidTag(t_name) {
		name = t_name
	}

	return []byte(name), mandatory
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
