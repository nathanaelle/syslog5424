package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"fmt"
	"testing"
)

type PriorityTest struct {
	a string
	p Priority
}

func TestPriority(t *testing.T) {
	lInval := []string{
		"foo.bar",
		"kern.emerg.info",
		"kern",
	}

	lVal := []PriorityTest{
		{"kern.emerg", Priority(0)},
		{"user.debug", Priority(15)},
	}

	d := new(Priority)

	for _, inv := range lInval {
		err := d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _, val := range lVal {
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

func TestPriorityMarshal5424(t *testing.T) {
	i := int(0)

	for i < 256 {
		z := Priority(i)
		m, err := z.Marshal5424()
		if err != nil {
			t.Errorf("%d marshal got err: %v", i, err)
			return
		}
		a := string(m)
		b := fmt.Sprintf("<%d>1", i)
		if a != b {
			t.Errorf("got [%v] expected [%v]", a, b)
		}
		i++
	}
}

/*

go test -cpu 4 -benchtime=10s -bench=Priority_ -benchmem
goarch: amd64
BenchmarkPrioritySet-4             	200000000	        66.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPriorityString-4          	300000000	        39.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPriorityMarshal5424-4     	500000000	        30.3 ns/op	       5 B/op	       1 allocs/op
BenchmarkPriorityUnmarshal5424-4   	2000000000	         9.86 ns/op	       0 B/op	       0 allocs/op
PASS

go test -cpu 1 -benchtime=10s -bench=Priority_ -benchmem
goarch: amd64
BenchmarkPrioritySet           		200000000	        67.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkPriorityString        		500000000	        33.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkPriorityMarshal5424   		500000000	        27.3 ns/op	       5 B/op	       1 allocs/op
BenchmarkPriorityUnmarshal5424 		2000000000	         8.47 ns/op	       0 B/op	       0 allocs/op
PASS

*/

func BenchmarkPrioritySet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := Priority(0)
		if err := (&p).Set("cron.warning"); err != nil {
			panic(err)
		}
	}
}

func BenchmarkPriorityString(b *testing.B) {
	p := Priority(15)
	for i := 0; i < b.N; i++ {
		if s := p.String(); s != "user.debug" {
			panic("benchmark expect user.debug")
		}
	}
}

func BenchmarkPriorityMarshal5424(b *testing.B) {
	p := Priority(15)
	for i := 0; i < b.N; i++ {
		if b, err := p.Marshal5424(); err != nil || string(b) != "<15>1" {
			panic("benchmark expect valid marshal")
		}
	}
}

func BenchmarkPriorityUnmarshal5424(b *testing.B) {
	data := []byte("<15>1")
	for i := 0; i < b.N; i++ {
		p := Priority(0)
		if err := p.Unmarshal5424(data); err != nil || p != Priority(15) {
			panic("benchmark expect valid unmarshal")
		}
	}
}
