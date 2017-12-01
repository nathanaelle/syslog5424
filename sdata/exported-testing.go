package sdata // import "github.com/nathanaelle/syslog5424/sdata"

import (
	"fmt"
	"reflect"
)

type (
	Foundable interface {
		String() string
		Found([]byte) (StructuredData, bool)
	}

	ValidTest struct {
		Orig string
		Obj  StructuredData
		Enc  string
	}

	SDTest interface {
		DoTest(sdid Foundable) error
	}

	InvalidMarshal struct {
		Obj StructuredData
		Err error
	}

	InvalidUnmarshal struct {
		Orig string
	}

	genericTest struct {
		Orig string
		Obj  StructuredData
	}
)

func (val ValidTest) DoTest(sdid Foundable) error {
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

func (val InvalidMarshal) DoTest(sdid Foundable) error {
	_, err := val.Obj.Marshal5424()
	if err != val.Err {
		return fmt.Errorf("expect err : %v got : %v", val.Err, err)
	}
	return nil
}

func (val InvalidUnmarshal) DoTest(sdid Foundable) error {
	_, ok := sdid.Found([]byte(val.Orig))
	if ok {
		return fmt.Errorf("%v SHOULDNT PARSE : %v", sdid, val.Orig)
	}
	return nil
}
