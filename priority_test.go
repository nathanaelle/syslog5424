package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"testing"
)

type PriorityTest struct {
	a string
	p Priority
}

func Test_Priority(t *testing.T) {
	l_inval := []string{
		"foo.bar",
		"kern.emerg.info",
		"kern",
	}

	l_val := []PriorityTest{
		PriorityTest{"kern.emerg", Priority(0)},
		PriorityTest{"user.debug", Priority(15)},
	}

	d := new(Priority)

	for _, inv := range l_inval {
		err := d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _, val := range l_val {
		d := new(Priority)
		err := d.Set(val.a)
		if err != nil {
			t.Errorf("[%v] parser invalid", val.a)
		}

		if val.a != d.String() {
			t.Errorf("[%v] [%v] differs", val.a, d)
		}
	}
}
