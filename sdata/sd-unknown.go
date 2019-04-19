package sdata // import "github.com/nathanaelle/syslog5424/v2/sdata"

type (
	unknownDef struct{}

	unknownSD struct {
		Name string
		Map  []struct {
			K string
			V string
		}
	}
)

var unknownSDID SDID = unknownDef{}

func (d unknownDef) String() string {
	return "unknown Structured Data"
}

func (d unknownDef) Options() map[string]interface{} {
	return map[string]interface{}{}
}

func (d unknownDef) SetOptions(_ map[string]interface{}) SDID {
	return d
}

func (d unknownDef) Default() unknownSD {
	return unknownSD{"", nil}
}

func (d unknownDef) IsIANA() bool {
	return false
}

func (d unknownDef) GetPEN() uint64 {
	return 32473
}

func (d unknownDef) MarshalText() (text []byte, err error) {
	return []byte("unknown@32473"), nil
}

func (d unknownDef) Found(data []byte) (StructuredData, bool) {
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

	ret := unknownSD{header, nil}

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

		ret.Map = append(ret.Map, struct{ K, V string }{name, value})
	}

	return ret, true
}

func (d unknownSD) SDID() SDID {
	return unknownSDID
}

func (d unknownSD) Marshal5424() ([]byte, error) {
	ret := "[" + d.Name

	for _, pair := range d.Map {
		ret += " " + pair.K + "=\"" + pair.V + "\""
	}

	ret += "]"

	return []byte(ret), nil
}
