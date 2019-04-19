package timequality

import (
	"testing"

	"github.com/nathanaelle/syslog5424/v2/sdata"
)

func pint(a int) *int {
	return &a
}

var sdt = []sdata.SDTest{
	sdata.TestValid(`[timeQuality]`, TimeQuality{false, false, nil}, `[timeQuality]`),
	sdata.TestValid(`[timeQuality tzKnown="0" isSynced="0"]`, TimeQuality{false, false, nil}, `[timeQuality]`),
	sdata.TestValid(`[timeQuality tzKnown="0"]`, TimeQuality{false, false, nil}, `[timeQuality]`),
	sdata.TestValid(`[timeQuality tzKnown="1" isSynced="0"]`, TimeQuality{true, false, nil}, `[timeQuality tzKnown="1"]`),
	sdata.TestValid(`[timeQuality tzKnown="1"]`, TimeQuality{true, false, nil}, `[timeQuality tzKnown="1"]`),
	sdata.TestValid(`[timeQuality tzKnown="1" isSynced="1"]`, TimeQuality{true, true, nil}, `[timeQuality tzKnown="1" isSynced="1"]`),
	sdata.TestValid(
		`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="60000000"]`,
		TimeQuality{true, true, pint(60000000)},
		`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="60000000"]`,
	),
	sdata.TestInvalidUnmarshal(`timequality]`),
	sdata.TestInvalidUnmarshal(`[timequality]`),
	sdata.TestInvalidUnmarshal(`[timequality`),
	sdata.TestInvalidUnmarshal(`[ timeQuality]`),
	sdata.TestInvalidUnmarshal(`[timeQuality invalid="1"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality invalid=1]`),
	sdata.TestInvalidUnmarshal(`[timeQuality invalid="1]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown ="1"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown="a"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown="1" isSynced="a"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown="0" isSynced="1"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown="1"isSynced="1"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="a"]`),
	sdata.TestInvalidUnmarshal(`[timeQuality tzKnown="1" isSynced="0" syncAccuracy="1000"]`),
	sdata.TestInvalidMarshal(TimeQuality{false, true, nil}, ErrInvalidTimeQuality),
	sdata.TestInvalidMarshal(TimeQuality{false, true, pint(100)}, ErrInvalidTimeQuality),
	sdata.TestInvalidMarshal(TimeQuality{false, false, pint(100)}, ErrInvalidTimeQuality),
	sdata.TestInvalidMarshal(TimeQuality{true, false, pint(100)}, ErrInvalidTimeQuality),
}

func TestTimeQuality(t *testing.T) {
	for i, val := range sdt {
		if err := val.DoTest(TQ); err != nil {
			t.Errorf("test %3d failed : %v", i, err)
		}
	}
}
