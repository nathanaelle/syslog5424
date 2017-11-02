package sdata // import "github.com/nathanaelle/syslog5424/sdata"

type (
	SDIDLight interface {
		String() string
		IsIANA() bool
		GetPEN() uint64
	}

	genericSDID struct {
		SDIDLight
	}

	genericSD struct {
		sdid      genericSDID
		Interface interface{}
	}
)

func GenericSD(i SDIDLight) StructuredData {
	return genericSD{genericSDID{i}, i}
}

func (sdid genericSDID) String() string {
	return sdid.SDIDLight.String()
}

func (sdid genericSDID) Options() map[string]interface{} {
	return map[string]interface{}{}
}

func (sdid genericSDID) SetOptions(_ map[string]interface{}) SDID {
	return sdid
}

func (sdid genericSDID) Default() StructuredData {
	return genericSD{sdid, nil}
}

func (sdid genericSDID) IsIANA() bool {
	return sdid.SDIDLight.IsIANA()
}

func (sdid genericSDID) GetPEN() uint64 {
	return sdid.SDIDLight.GetPEN()
}

func (sdid genericSDID) MarshalText() (text []byte, err error) {
	return []byte(sdid.String()), nil
}

func (sdid genericSDID) Found(data []byte) (StructuredData, bool) {
	if data[0] != '[' || data[len(data)-1] != ']' {
		return nil, false
	}
	data = data[1 : len(data)-1]

	header := sdid.String()
	if string(data[0:len(header)]) != header {
		return nil, false
	}

	data = data[len(header):]
	if len(data) == 0 {
		return genericSD{sdid, nil}, true
	}

	return nil, false
}

func (sd genericSD) SDID() SDID {
	return sd.sdid
}

func (sd genericSD) Marshal5424() ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(sd.Interface)
	if err != nil {
		return []byte{}, err
	}
	return e.Bytes(), nil
}
