package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"math"
	"testing"
	"time"
)

type MessageTest struct {
	m Message
	a string
}

var messageTest = []MessageTest{
	{
		Message{Priority(0), z_epoch(), "-", "-", "-", "-", emptyListSD, ""},
		"<0>1 1970-01-01T01:00:00Z - - - - -",
	},
	{
		Message{Priority(0), z_epoch(), "-", "-", "-", "-", []interface{}{timeQuality{pint(1), pint(1), nil}}, ""},
		`<0>1 1970-01-01T01:00:00Z - - - - [timeQuality tzKnown="1" isSynced="1"]`,
	},
	{
		Message{Priority(24), z_epoch(), "bla", "bli", "blu", "blo", emptyListSD, "message"},
		"<24>1 1970-01-01T01:00:00Z bla bli blu blo - message",
	},
}

var parseTest = []string{
	"<0>1 1970-01-01T01:00:00+01:00 - - - - -",
	"<12>1 1970-01-01T01:00:00Z - - - - -",
	"<0>1 1970-01-01T01:00:00Z bla bli blu blo - message",
	"<234>1 1970-01-01T01:00:00+01:00 bla bli blu blo - message",
}

func z_epoch() time.Time {
	t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "1970-01-01T01:00:00Z", time.UTC)
	return t
}

func Test_Message(t *testing.T) {
	for _, tt := range messageTest {
		a := tt.m.String()
		if a != tt.a {
			t.Errorf(" %v String() = %s; want %s", tt.m, a, tt.a)
			continue
		}
	}

	for _, tt := range parseTest {
		a, err := Parse([]byte(tt))
		if err != nil {
			t.Errorf(" %s parse() %s", tt, err)
		}

		if a.String() != tt {
			t.Errorf(" %s parse() %+v [%s]", tt, a, a.String())
			continue
		}
	}
}

func Benchmark_Message_Parse(b *testing.B) {
	max := int(math.Ceil(float64(b.N) / float64(len(parseTest))))
	for i := 0; i < max; i++ {
		for _, tt := range parseTest {
			Parse([]byte(tt))
		}
	}
}

func Benchmark_Message_Chan(b *testing.B) {
	max := int(math.Ceil(float64(b.N) / float64(len(messageTest))))
	ch := make(chan Message)
	go func() {
		for {
			<-ch
		}
	}()
	for i := 0; i < max; i++ {
		for _, tt := range messageTest {
			ch <- tt.m
		}
	}
}

func Benchmark_Message_ChanBuf(b *testing.B) {
	max := int(math.Ceil(float64(b.N) / float64(len(messageTest))))
	ch := make(chan Message, 100)
	go func() {
		for {
			<-ch
		}
	}()
	for i := 0; i < max; i++ {
		for _, tt := range messageTest {
			ch <- tt.m
		}
	}
}

func Benchmark_Message_String(b *testing.B) {
	max := int(math.Ceil(float64(b.N) / float64(len(messageTest))))
	for i := 0; i < max; i++ {
		for _, tt := range messageTest {
			tt.m.String()
		}
	}
}

func Benchmark_Message_String_sd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messageTest[1].m.String()
	}
}

func Benchmark_Message_String_short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messageTest[0].m.String()
	}
}

func Benchmark_Message_String_long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messageTest[2].m.String()
	}
}

func Benchmark_Message_CreateMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateMessage("test", Priority(0), "test")
	}
}

func Benchmark_Message_CreateMessage_Call(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EmptyMessage().AppName("test-app").Priority(Priority(21)).LocalHost().Now().Msg("test message")
	}
}
