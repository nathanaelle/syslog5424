package sdata // import "github.com/nathanaelle/syslog5424/sdata"

import ()

type (
	unknownDef struct{}

	unknownSD struct {
		Name string
		Map  map[string][]string
	}
)

var unknown_sdid SDID = unknownDef{}

func (_ unknownDef) String() string {
	return "unknown Structured Data"
}

func (_ unknownDef) Options() map[string]interface{} {
	return map[string]interface{}{}
}

func (d unknownDef) SetOptions(_ map[string]interface{}) SDID {
	return d
}

func (d unknownDef) Default() unknownSD {
	return unknownSD{"", nil}
}

func (_ unknownDef) IsIANA() bool {
	return false
}

func (_ unknownDef) GetPEN() uint64 {
	return 32473
}

func (_ unknownDef) MarshalText() (text []byte, err error) {
	return []byte("unknown@32473"), nil
}

func (_ unknownDef) Found(data []byte) (StructuredData, bool) {
	if data[0] != '[' || data[len(data)-1] != ']' {
		return nil, false
	}
	data = data[1 : len(data)-1]

	header, data, err := NextHeader(data)
	if err != nil {
		return nil, false
	}

	if len(data) == 0 {
		return unknownSD{header, nil}, true
	}

	ret := unknownSD{header, make(map[string][]string)}

	for len(data) > 0 {
		var err error
		data, err = NextNonSpace(data)
		if err != nil {
			return nil, false
		}

		name, step1, err := NextSDName(data)
		if err != nil {
			return nil, false
		}
		value, step2, err := NextSDValue(step1)
		if err != nil {
			return nil, false
		}
		data = step2

		ret.Map[name] = append(ret.Map[name], value)
	}

	return ret, true
}

func (_ unknownSD) SDID() SDID {
	return unknown_sdid
}

func (sd unknownSD) Marshal5424() ([]byte, error) {
	ret := "[" + sd.Name

	for key, values := range sd.Map {
		for _, val := range values {
			ret += " " + key + "=\"" + val + "\""
		}
	}

	ret += "]"

	return []byte(ret), nil
}
