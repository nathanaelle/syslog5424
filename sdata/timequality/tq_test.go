package timequality // import "github.com/nathanaelle/syslog5424/sdata/timequality"

import (
	"testing"
	//	"github.com/nathanaelle/syslog5424/sdata"
	"../../sdata"
)

func pint(a int) *int {
	return &a
}

var sd_t = []sdata.SDTest{
	sdata.ValidTest{
		`[timeQuality]`, TimeQuality{false, false, nil}, `[timeQuality]`,
	},
	sdata.ValidTest{
		`[timeQuality tzKnown="0" isSynced="0"]`, TimeQuality{false, false, nil}, `[timeQuality]`,
	},
	sdata.ValidTest{
		`[timeQuality tzKnown="0"]`, TimeQuality{false, false, nil}, `[timeQuality]`,
	},
	sdata.ValidTest{
		`[timeQuality tzKnown="1" isSynced="0"]`, TimeQuality{true, false, nil}, `[timeQuality tzKnown="1"]`,
	},
	sdata.ValidTest{
		`[timeQuality tzKnown="1"]`, TimeQuality{true, false, nil}, `[timeQuality tzKnown="1"]`,
	},
	sdata.ValidTest{
		`[timeQuality tzKnown="1" isSynced="1"]`, TimeQuality{true, true, nil}, `[timeQuality tzKnown="1" isSynced="1"]`,
	},
	sdata.ValidTest{
		`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="60000000"]`,
		TimeQuality{true, true, pint(60000000)},
		`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="60000000"]`,
	},
	sdata.InvalidUnmarshal{`timequality]`},
	sdata.InvalidUnmarshal{`[timequality]`},
	sdata.InvalidUnmarshal{`[timequality`},
	sdata.InvalidUnmarshal{`[ timeQuality]`},
	sdata.InvalidUnmarshal{`[timeQuality invalid="1"]`},
	sdata.InvalidUnmarshal{`[timeQuality invalid=1]`},
	sdata.InvalidUnmarshal{`[timeQuality invalid="1]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown ="1"]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown="a"]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown="1" isSynced="a"]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown="0" isSynced="1"]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown="1"isSynced="1"]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="a"]`},
	sdata.InvalidUnmarshal{`[timeQuality tzKnown="1" isSynced="0" syncAccuracy="1000"]`},
	sdata.InvalidMarshal{TimeQuality{false, true, nil}, InvalidTimeQuality},
	sdata.InvalidMarshal{TimeQuality{false, true, pint(100)}, InvalidTimeQuality},
	sdata.InvalidMarshal{TimeQuality{false, false, pint(100)}, InvalidTimeQuality},
	sdata.InvalidMarshal{TimeQuality{true, false, pint(100)}, InvalidTimeQuality},
}

func TestTimeQuality(t *testing.T) {
	for i, val := range sd_t {
		if err := val.DoTest(TQ); err != nil {
			t.Errorf("test %3d failed : %v", i, err)
		}
	}
}
