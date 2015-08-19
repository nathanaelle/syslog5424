package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"testing"
)

type origin map[string][]string

type timeQuality struct {
	TzKnown      *int `sd5424:"tzKnown"`
	IsSynced     *int `sd5424:"isSynced"`
	SyncAccuracy *int `sd5424:"syncAccuracy"`
}

type exampleSDID map[string]string

func (exampleSDID) GetPEN() uint64 {
	return uint64(32473)
}

func (exampleSDID) SetPEN(_ uint64) {
}

func pint(a int) *int {
	p := new(int)
	*p = a
	return p
}

type SDTest struct {
	asc string
	obj interface{}
}

var l_val []SDTest = []SDTest{
	SDTest{
		`[exampleSDID@32473 eventID="1011" eventSource="Application" iut="3"]`,
		exampleSDID{"iut": "3", "eventSource": "Application", "eventID": "1011"},
	},
	SDTest{
		`[exampleSDID@32473 eventID="1011" eventSource="[Application\]" iut="3"]`,
		exampleSDID{"iut": "3", "eventSource": "[Application]", "eventID": "1011"},
	},
	SDTest{
		`[timeQuality]`,
		timeQuality{},
	},
	SDTest{
		`[timeQuality tzKnown="0" isSynced="0"]`,
		timeQuality{pint(0), pint(0), nil},
	},
	SDTest{
		`[timeQuality tzKnown="1" isSynced="1"]`,
		timeQuality{pint(1), pint(1), nil},
	},
	SDTest{
		`[timeQuality tzKnown="1" isSynced="1" syncAccuracy="60000000"]`,
		timeQuality{pint(1), pint(1), pint(60000000)},
	},
	SDTest{
		`[origin ip="192.0.2.1" ip="192.0.2.129"]`,
		origin{"ip": []string{"192.0.2.1", "192.0.2.129"}},
	},
}

func Test_SDMarshal(t *testing.T) {
	for _, val := range l_val {
		d := string(MarshalSD(val.obj))
		if d != val.asc {
			t.Errorf("%v [%v] differs", val.asc, d)
		}
	}
}

func Test_SDUnmarshal(t *testing.T) {
	t.Skip()
	/*	l_inval	:= []string{
			`[ exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"]`,
			`[ exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"]`,
			`[exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] [examplePriority@32473 class="high"]`,
		}

		for _,inv := range l_inval {
			err	:= d.Set(inv)
			if err == nil {
				t.Errorf("[%v] parser invalid", inv)
			}
		}
	*/
}
