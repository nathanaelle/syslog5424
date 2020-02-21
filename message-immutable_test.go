package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"testing"
	"time"

	"github.com/nathanaelle/syslog5424/v2/sdata/timequality"
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
	`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][exampleEscape@32473 text="some[data\] here"]`,
}

func Test_MessageImmutable_Parse(t *testing.T) {
	for _, tt := range parseTest {
		a, _, err := Parse([]byte(tt), nil, false)
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

func Test_MessageImmutable_Values(t *testing.T) {
	tt := `<12>1 1970-01-01T01:00:00Z - - - - -`
	a, _, err := Parse([]byte(tt), nil, false)
	if err != nil {
		t.Errorf("parse() {{ %q }} is : %q", tt, err)
		return
	}

	if a.String() != tt {
		t.Errorf(" %q parse() %+v [%q]", tt, a, a.String())
		return
	}

	msg := a.Writable().String()
	if msg != tt {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", tt, msg)
		return
	}

	if a.Priority() != (LogUSER | LogWARNING) {
		t.Errorf("Writable: expected {{ %s }} got {{ %s }}", (LogUSER | LogWARNING), a.Priority())
		return
	}

	if !a.TimeStamp().Equal(time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)) {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC), a.TimeStamp())
		return
	}

	if a.Hostname() != "-" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "-", a.Hostname())
		return
	}

	if a.AppName() != "-" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "-", a.AppName())
		return
	}

	if a.ProcID() != "-" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "-", a.ProcID())
		return
	}

	if a.MsgID() != "-" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "-", a.MsgID())
		return
	}

	if a.StructuredData().String() != "-" {
		t.Errorf("Writable: expected {{ %q }} got {{ %+v }}", "-", a.StructuredData().String())
		return
	}

	if a.Text() != "" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "", a.Text())
		return
	}

	tt = `<165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 myproc 8710 foobar [timeQuality tzKnown="1" isSynced="1"] %% It's time to make the do-nuts.`
	a, _, err = Parse([]byte(tt), nil, false)
	if err != nil {
		t.Errorf("parse() {{ %q }} is : %q", tt, err)
		return
	}

	if a.String() != tt {
		t.Errorf(" %q parse() %+v [%q]", tt, a, a.String())
		return
	}

	msg = a.Writable().String()
	if msg != tt {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", tt, msg)
		return
	}

	if a.Priority() != (LogLOCAL4 | LogNOTICE) {
		t.Errorf("Writable: expected {{ %s }} got {{ %s }}", (LogLOCAL4 | LogNOTICE), a.Priority())
		return
	}

	if !a.TimeStamp().UTC().Equal(time.Date(2003, 8, 24, 12, 14, 15, 3000, time.UTC)) {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", time.Date(2003, 8, 24, 12, 14, 15, 3000, time.UTC), a.TimeStamp().UTC())
		return
	}

	if a.Hostname() != "192.0.2.1" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "192.0.2.1", a.Hostname())
		return
	}

	if a.AppName() != "myproc" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "myproc", a.AppName())
		return
	}

	if a.ProcID() != "8710" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "8710", a.ProcID())
		return
	}

	if a.MsgID() != "foobar" {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "foobar", a.MsgID())
		return
	}

	msgsd := a.StructuredData()
	if len(msgsd) != 1 {
		t.Errorf("Writable: expected {{ %q }} got {{ %+v }}", "-", a.StructuredData().String())
		return
	}

	if _, ok := msgsd[0].(timequality.TimeQuality); !ok {
		t.Errorf("Writable: expected {{ %q }} got {{ %+v }}", "-", a.StructuredData().String())
		return
	}

	if a.Text() != "%% It's time to make the do-nuts." {
		t.Errorf("Writable: expected {{ %q }} got {{ %q }}", "%% It's time to make the do-nuts.", a.Text())
		return
	}
}

/*

Actual performance are :

go test -cpu 1 -benchtime=10s -bench=MessageImmutable_Parse -benchmem
goarch: amd64
pkg: github.com/nathanaelle/syslog5424/v2
Benchmark_MessageImmutable_Parse_Minimal                        	100000000	       117 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MessageNoSD                    	89648371	       138 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_KnownSDOnly                    	87056971	       140 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_UnkownSDOnly                   	88181227	       140 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MessageAndSD                   	75508340	       163 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_Verylong                       	59195804	       188 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MaxLength                      	64618474	       176 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MessageAndSD_Then_Make_Mutable 	 5349686	      2237 ns/op	     736 B/op	      24 allocs/op
PASS

go test -cpu 4 -benchtime=10s -bench=MessageImmutable_Parse -benchmem
goarch: amd64
pkg: github.com/nathanaelle/syslog5424/v2
Benchmark_MessageImmutable_Parse_Minimal-4                          	100000000	       113 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MessageNoSD-4                      	91566769	       130 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_KnownSDOnly-4                      	90291741	       130 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_UnkownSDOnly-4                     	92336280	       129 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MessageAndSD-4                     	76834656	       154 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_Verylong-4                         	67873761	       176 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MaxLength-4                        	72529458	       168 ns/op	      64 B/op	       1 allocs/op
Benchmark_MessageImmutable_Parse_MessageAndSD_Then_Make_Mutable-4   	 5741421	      2106 ns/op	     736 B/op	      24 allocs/op
PASS

*/

func Benchmark_MessageImmutable_Parse_Minimal(b *testing.B) {
	var err error

	data := []byte(`<0>1 1970-01-01T01:00:00+01:00 - - - - -`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_MessageNoSD(b *testing.B) {
	var err error

	data := []byte(`<165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 myproc 8710 - - %% It's time to make the do-nuts.`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_KnownSDOnly(b *testing.B) {
	var err error

	data := []byte(`<0>1 1970-01-01T01:00:00Z - - - - [timeQuality tzKnown="1" isSynced="1"]`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_UnkownSDOnly(b *testing.B) {
	var err error

	data := []byte(`<0>1 1970-01-01T01:00:00Z - - - - [timeQualitat tzKnown="1" isSynced="1"]`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_MessageAndSD(b *testing.B) {
	var err error

	data := []byte(`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"] Some log message with structured data`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_Verylong(b *testing.B) {
	var err error

	data := []byte(`<190>1 2016-02-21T01:19:11+00:00 batch6sj - - - [meta sequenceId="21881798" x-group="37051387"][origin x-service="tracking"] metascutellar conversationalist nephralgic exogenetic graphy streng outtaken acouasm amateurism prenotice Lyonese bedull antigrammatical diosphenol gastriloquial bayoneteer sweetener naggy roughhouser dighter addend sulphacid uneffectless ferroprussiate reveal Mazdaist plaudite Australasian distributival wiseman rumness Seidel topazine shahdom sinsion mesmerically pinguedinous ophthalmotonometer scuppler wound eciliate expectedly carriwitchet dictatorialism bindweb pyelitic idic atule kokoon poultryproof rusticial seedlip nitrosate splenadenoma holobenthic uneternal Phocaean epigenic doubtlessly indirection torticollar robomb adoptedly outspeak wappenschawing talalgia Goop domitic savola unstrafed carded unmagnified mythologically orchester obliteration imperialine undisobeyed galvanoplastical cycloplegia quinquennia foremean umbonal marcgraviaceous happenstance theoretical necropoles wayworn Igbira pseudoangelic raising unfrounced lamasary centaurial Japanolatry microlepidoptera`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_MaxLength(b *testing.B) {
	var err error

	data := []byte(`<191>1 2018-12-31T23:59:59.999999-23:59 abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabc abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdef abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzab abcdefghilmnopqrstuvzabcdefghilm [an@id key1="val1" key2="val2"][another@id key1="val1"] Some message "GET"`)
	for i := 0; i < b.N; i++ {
		_, _, err = Parse(data, nil, false)
	}
	if err != nil {
		panic(err)
	}
}

func Benchmark_MessageImmutable_Parse_MessageAndSD_Then_Make_Mutable(b *testing.B) {
	data := []byte(`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"] Some log message with structured data`)
	for i := 0; i < b.N; i++ {
		mi, _, err := Parse(data, nil, false)
		if err != nil {
			panic(err)
		}
		mi.Writable()
	}
}
