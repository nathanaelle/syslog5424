package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Message struct {
	prio      Priority
	timestamp time.Time
	hostname  string
	appname   string
	procid    string
	msgid     string
	SD        listStructuredData
	message   string
}

const RFC5424TimeStamp string = "2006-01-02T15:04:05.999999Z07:00"

var hostname, _ = os.Hostname()

func EmptyMessage() Message {
	return Message{Priority(0), time.Unix(0, 0), "-", "-", "-", "-", emptyListSD, ""}
}

func (msg Message) Now() Message {
	return Message{msg.prio, time.Now(), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func stamp_to_ts(stamp string) time.Time {
	now := time.Now()
	ts, _ := time.Parse(time.Stamp, stamp)
	year := now.Year()

	if now.Month() == 1 && ts.Month() == 12 {
		year--
	}

	return time.Date(year, ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), ts.Location())
}

func (msg Message) Stamp(stamp string) Message {
	return Message{msg.prio, stamp_to_ts(stamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func delta_boot_to_ts(boot_ts time.Time, s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return boot_ts.Add(time.Duration(nsec)*time.Nanosecond + time.Duration(sec)*time.Second)
}

func (msg Message) Delta(boot_ts time.Time, s_sec string, s_nsec string) Message {
	return Message{msg.prio, delta_boot_to_ts(boot_ts, s_sec, s_nsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func epoc_to_ts(s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return time.Unix(sec, nsec)
}

func (msg Message) Epoch(s_sec string, s_nsec string) Message {
	return Message{msg.prio, epoc_to_ts(s_sec, s_nsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func (msg Message) App(appname string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func (msg Message) ProcID(procid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, procid, msg.msgid, msg.SD, msg.message}
}

func (msg Message) MsgID(msgid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msgid, msg.SD, msg.message}
}

func (msg Message) Priority(prio Priority) Message {
	return Message{prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func (msg Message) LocalHost() Message {
	return Message{msg.prio, msg.timestamp, hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message}
}

func (msg Message) Data(data string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, data}
}

func (msg Message) String() string {
	switch msg.message {
	case "":
		return fmt.Sprintf("<%v>1 %s %s %s %s %s %v",
			msg.prio, msg.timestamp.Format(RFC5424TimeStamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD)

	default:
		return fmt.Sprintf("<%v>1 %s %s %s %s %s %v %s",
			msg.prio, msg.timestamp.Format(RFC5424TimeStamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.SD, msg.message)

	}

}

func CreateMessage(data string, appname string, prio Priority) Message {
	return EmptyMessage().App(appname).Priority(prio).LocalHost().Now().Data(data)
}
