package sdata // import "github.com/nathanaelle/syslog5424/v2/sdata"

import (
	"fmt"
	"reflect"
)

type (
	// Foundable is the minimal part of the SDID interface for testing
	Foundable interface {
		String() string
		Found([]byte) (StructuredData, bool)
	}

	validTest struct {
		Orig string
		Obj  StructuredData
		Enc  string
	}

	// SDTest describe the common interface of any test on Structured Data
	SDTest interface {
		DoTest(sdid Foundable) error
	}

	invalidMarshal struct {
		Obj StructuredData
		Err error
	}

	invalidUnmarshal struct {
		Orig string
	}

	genericTest struct {
		Orig string
		Obj  StructuredData
	}
)

// TestValid define a test for valid Structured Data
func TestValid(Orig string, Obj StructuredData, Enc string) SDTest {
	return validTest{
		Orig: Orig,
		Obj:  Obj,
		Enc:  Enc,
	}
}

// TestInvalidUnmarshal define a test a non Unmarshallable Structured Data
func TestInvalidUnmarshal(Orig string) SDTest {
	return invalidUnmarshal{
		Orig: Orig,
	}
}

// TestInvalidMarshal define a test a non Marshallable Structured Data
func TestInvalidMarshal(Obj StructuredData, Err error) SDTest {
	return invalidMarshal{
		Obj: Obj,
		Err: Err,
	}
}

func (val validTest) DoTest(sdid Foundable) error {
	o, ok := sdid.Found([]byte(val.Orig))
	if !ok {
		return fmt.Errorf("%v can't parse : %v", sdid, val.Orig)
	}

	if !reflect.DeepEqual(o, val.Obj) {
		return fmt.Errorf("%#v and %#v differs", o, val.Obj)
	}

	d, err := o.Marshal5424()
	if err != nil {
		return fmt.Errorf("%v got err : %v", val.Orig, err)
	}

	if string(d) != val.Enc {
		return fmt.Errorf("%v and %v differs", val.Enc, string(d))
	}

	return nil
}

func (val genericTest) DoTest(sdid Foundable) error {
	/*	o, ok := sdid.Found([]byte(val.Orig))
		if !ok {
			return	fmt.Errorf("%v can't parse : %v", sdid, val.Orig)
		}

		if !reflect.DeepEqual(o, val.Obj) {
			return	fmt.Errorf("%#v and %#v differs", o, val.Obj)
		}
	*/
	d, err := val.Obj.Marshal5424()
	if err != nil {
		return fmt.Errorf("%v got err : %v", val.Orig, err)
	}

	if string(d) != val.Orig {
		return fmt.Errorf("%v and %v differs", val.Orig, string(d))
	}

	return nil
}

func (val invalidMarshal) DoTest(sdid Foundable) error {
	_, err := val.Obj.Marshal5424()
	if err != val.Err {
		return fmt.Errorf("expect err : %v got : %v", val.Err, err)
	}
	return nil
}

func (val invalidUnmarshal) DoTest(sdid Foundable) error {
	_, ok := sdid.Found([]byte(val.Orig))
	if ok {
		return fmt.Errorf("%v SHOULDNT PARSE : %v", sdid, val.Orig)
	}
	return nil
}
