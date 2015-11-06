package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"fmt"
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

func Test_Priority_Marshal5424(t *testing.T) {
	i	:= int(0)

	for i < 256 {
		z :=Priority(i)
		a := string(z.Marshal5424())
		b := fmt.Sprintf("<%d>1", i)
		if a != b {
			t.Errorf("m[%v] f[%v] differs", a, b)
		}
		i++
	}
}
