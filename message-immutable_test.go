package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	_ "./sdata/timequality"
	"testing"
)

var parseTest = []string{
	`<0>1 1970-01-01T01:00:00+01:00 - - - - -`,
	`<12>1 1970-01-01T01:00:00Z - - - - -`,
	`<0>1 1970-01-01T01:00:00Z bla bli blu blo - message`,
	`<234>1 1970-01-01T01:00:00+01:00 bla bli blu blo - message`,
	`<0>1 1970-01-01T01:00:00Z - - - - [timeQuality tzKnown="1" isSynced="1"]`,
	`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - \xEF\xBB\xBF'su root' failed for lonvick on /dev/pts/8`,
	`<165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 myproc 8710 - - %% It's time to make the do-nuts.`,
	`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] \xEF\xBB\xBFAn application event log entry...`,
	`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"]`,
	`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][exampleEscape@32473 text="some[data\]here"]`,
}

func TestMessageImmutableParse(t *testing.T) {
	for _, tt := range parseTest {
		a, err := Parse([]byte(tt))
		if err != nil {
			t.Errorf("parse() {{ %q }} is : %q", tt, err)
			continue
		}

		if a.String() != tt {
			t.Errorf(" %q parse() %+v [%q]", tt, a, a.String())
			continue
		}

		msg := a.Writable().String()
		if msg != tt {
			t.Errorf("Writable: expected {{ %q }} got {{ %q }}", tt, msg)
			continue
		}
	}
}
