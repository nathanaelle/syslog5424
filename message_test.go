package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"math"
	"testing"
	"time"

	"github.com/nathanaelle/syslog5424/v2/sdata"
	tq "github.com/nathanaelle/syslog5424/v2/sdata/timequality"
)

type MessageTest struct {
	m Message
	a string
}

var messageTest = []MessageTest{
	{
		Message{Priority(0), zEpoch(), "-", "-", "-", "-", emptyListSD(), ""},
		"<0>1 1970-01-01T01:00:00Z - - - - -",
	},
	{
		Message{Priority(0), zEpoch(), "-", "-", "-", "-", sdata.List{tq.TimeQuality{TzKnown: true, IsSynced: true, SyncAccuracy: nil}}, ""},
		`<0>1 1970-01-01T01:00:00Z - - - - [timeQuality tzKnown="1" isSynced="1"]`,
	},
	{
		Message{Priority(24), zEpoch(), "bla", "bli", "blu", "blo", emptyListSD(), "message"},
		"<24>1 1970-01-01T01:00:00Z bla bli blu blo - message",
	},
}

func zEpoch() time.Time {
	t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "1970-01-01T01:00:00Z", time.UTC)
	return t
}

func TestMessage(t *testing.T) {
	for _, tt := range messageTest {
		a := tt.m.String()
		if a != tt.a {
			t.Errorf(" %v String() = %s; want %s", tt.m, a, tt.a)
			continue
		}
	}
}

func Benchmark_Message_String(b *testing.B) {
	var str string

	max := int(math.Ceil(float64(b.N) / float64(len(messageTest))))
	for i := 0; i < max; i++ {
		for _, tt := range messageTest {
			str = tt.m.String()
		}
	}
	str += ""
}

func Benchmark_Message_String_sd(b *testing.B) {
	var str string

	for i := 0; i < b.N; i++ {
		str = messageTest[1].m.String()
	}
	str += ""
}

func Benchmark_Message_String_short(b *testing.B) {
	var str string

	for i := 0; i < b.N; i++ {
		str = messageTest[0].m.String()
	}
	str += ""
}

func Benchmark_Message_String_long(b *testing.B) {
	var str string

	for i := 0; i < b.N; i++ {
		str = messageTest[2].m.String()
	}
	str += ""
}

func Benchmark_Message_Marshal5424(b *testing.B) {
	msg := Message{Priority(0), zEpoch(), "-", "-", "-", "-", sdata.List{tq.TimeQuality{TzKnown: true, IsSynced: true, SyncAccuracy: nil}}, "It's time to make the do-nuts."}

	for i := 0; i < b.N; i++ {
		msg.Marshal5424()
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
