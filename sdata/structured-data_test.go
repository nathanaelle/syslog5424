package sdata // import "github.com/nathanaelle/syslog5424/sdata"

import (
	"testing"
)

type (
	origin map[string][]string

	exampleSDID map[string]string

	timeQuality map[string]string
)

func (exampleSDID) GetPEN() uint64 {
	return uint64(32473)
}

func (exampleSDID) IsIANA() bool {
	return false
}

func (exampleSDID) String() string {
	return "exampleSDID@32473"
}

func (timeQuality) GetPEN() uint64 {
	return uint64(0)
}

func (timeQuality) IsIANA() bool {
	return true
}

func (timeQuality) String() string {
	return "timeQuality"
}

func (origin) GetPEN() uint64 {
	return uint64(0)
}

func (origin) IsIANA() bool {
	return true
}

func (origin) String() string {
	return "origin"
}

func pint(a int) *int {
	return &a
}

var sd_t []SDTest = []SDTest{
	genericTest{
		`[exampleSDID@32473 eventID="1011" eventSource="Application" iut="3"]`,
		GenericSD(exampleSDID{"iut": "3", "eventSource": "Application", "eventID": "1011"}),
	},
	genericTest{
		`[exampleSDID@32473 eventID="1011" eventSource="[Application\]" iut="3"]`,
		GenericSD(exampleSDID{"iut": "3", "eventSource": "[Application]", "eventID": "1011"}),
	},
	genericTest{
		`[origin ip="192.0.2.1" ip="192.0.2.129"]`,
		GenericSD(origin{"ip": []string{"192.0.2.1", "192.0.2.129"}}),
	},
	InvalidUnmarshal{`[ exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"]`},
	InvalidUnmarshal{`[exampleSDID iut="3" eventSource ="Application" eventID="1011"]`},
}

func TestStructuredData(t *testing.T) {
	for i, val := range sd_t {
		if err := val.DoTest(registry); err != nil {
			t.Errorf("test %3d failed : %v", i, err)
		}
	}
}

func Benchmark_SDMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		/*for _, tt := range messageTest {
			tt.m.String()
		}*/
		//		MarshalSD(sd_t[3].obj)
	}
}
