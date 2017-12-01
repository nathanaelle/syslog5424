package timequality // import "github.com/nathanaelle/syslog5424/sdata/timequality"

import (
	"github.com/nathanaelle/syslog5424/sdata"
	"errors"
	"math"
	"strconv"
)

func init() {
	TQ = sdata.Register(tqdef{})
}

var TQ sdata.SDID
var InvalidTimeQuality error = errors.New("Invalid TimeQuality struct")

type (
	tqdef struct{}

	TimeQuality struct {
		TzKnown      bool `sd5424:"tzKnown"`
		IsSynced     bool `sd5424:"isSynced"`
		SyncAccuracy *int `sd5424:"syncAccuracy"`
	}
)

func (_ tqdef) String() string {
	return "TimeQuality"
}

func (_ tqdef) Options() map[string]interface{} {
	return map[string]interface{}{}
}

func (d tqdef) SetOptions(_ map[string]interface{}) sdata.SDID {
	return d
}

func (d tqdef) Default() TimeQuality {
	return TimeQuality{false, false, nil}
}

func (_ tqdef) IsIANA() bool {
	return true
}

func (_ tqdef) GetPEN() uint64 {
	return 0
}

func (_ tqdef) MarshalText() (text []byte, err error) {
	return []byte("timeQuality"), nil
}

func (_ tqdef) Found(data []byte) (sdata.StructuredData, bool) {
	if data[0] != '[' || data[len(data)-1] != ']' {
		return nil, false
	}
	data = data[1 : len(data)-1]

	if string(data[0:11]) != "timeQuality" {
		return nil, false
	}
	data = data[11:]
	if len(data) == 0 {
		return TimeQuality{false, false, nil}, true
	}

	tq := TimeQuality{false, false, nil}

	for len(data) > 0 {
		var err error
		data, err = sdata.NextNonSpace(data)
		if err != nil {
			return nil, false
		}

		name, step1, err := sdata.NextSDName(data)
		if err != nil {
			return nil, false
		}
		value, step2, err := sdata.NextSDValue(step1)
		if err != nil {
			return nil, false
		}
		data = step2

		switch name {
		case "tzKnown":
			switch value {
			case "0":
				tq.TzKnown = false
			case "1":
				tq.TzKnown = true
			default:
				return nil, false
			}
		case "isSynced":
			switch value {
			case "0":
				tq.IsSynced = false
			case "1":
				tq.IsSynced = true
			default:
				return nil, false
			}
		case "syncAccuracy":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, false
			}
			tq.SyncAccuracy = &v
		default:
			return nil, false
		}
	}

	if !tq.IsValid() {
		return nil, false
	}

	return tq, true
}

func (_ TimeQuality) SDID() sdata.SDID {
	return TQ
}

func (tq TimeQuality) IsValid() bool {
	return !((!tq.TzKnown && tq.IsSynced) || (!tq.IsSynced && tq.SyncAccuracy != nil))
}

func (tq TimeQuality) Marshal5424() ([]byte, error) {
	if !tq.IsValid() {
		return nil, InvalidTimeQuality
	}

	length := 13
	if tq.TzKnown {
		length += 12
	}
	if tq.IsSynced {
		length += 13
	}

	if tq.SyncAccuracy != nil {
		length += 15 + int(math.Floor(math.Log(float64(*tq.SyncAccuracy))+1))
	}

	ret := append(make([]byte, 0, length), []byte("[timeQuality")...)

	if tq.TzKnown {
		ret = append(ret, []byte(` tzKnown="1"`)...)
	}

	if tq.IsSynced {
		ret = append(ret, []byte(` isSynced="1"`)...)
	}

	if tq.SyncAccuracy != nil {
		ret = append(ret, []byte(` syncAccuracy="`)...)
		ret = strconv.AppendInt(ret, int64(*tq.SyncAccuracy), 10)
		ret = append(ret, '"')
	}
	ret = append(ret, ']')

	return ret, nil
}
