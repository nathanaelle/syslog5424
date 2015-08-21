package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
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
	sd        listStructuredData
	message   string
}

const RFC5424TimeStamp string = "2006-01-02T15:04:05.999999Z07:00"

var hostname, _ = os.Hostname()

func EmptyMessage() Message {
	return Message{Priority(0), time.Unix(0, 0), "-", "-", "-", "-", emptyListSD, ""}
}

func (msg Message) Now() Message {
	return Message{msg.prio, time.Now(), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
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
	return Message{msg.prio, stamp_to_ts(stamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

func delta_boot_to_ts(boot_ts time.Time, s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return boot_ts.Add(time.Duration(nsec)*time.Nanosecond + time.Duration(sec)*time.Second)
}

func (msg Message) Delta(boot_ts time.Time, s_sec string, s_nsec string) Message {
	return Message{msg.prio, delta_boot_to_ts(boot_ts, s_sec, s_nsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

func epoc_to_ts(s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return time.Unix(sec, nsec)
}

// set the date of a Message with a epoch TimeStamp
func (msg Message) Epoch(s_sec string, s_nsec string) Message {
	return Message{msg.prio, epoc_to_ts(s_sec, s_nsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// set the app-name of a Message
func (msg Message) AppName(appname string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// set the proc-id of a Message
func (msg Message) ProcID(procid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, procid, msg.msgid, msg.sd, msg.message}
}

// set the msg-id of a Message
func (msg Message) MsgID(msgid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msgid, msg.sd, msg.message}
}

// set the priority of a Message
func (msg Message) Priority(prio Priority) Message {
	return Message{prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

//set the hostname as the value get with gethostbyname()
func (msg Message) LocalHost() Message {
	return Message{msg.prio, msg.timestamp, hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

//set the message part of a Message
func (msg Message) Msg(message string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, message}
}

//set the message part of a Message
func (msg Message) StructuredData(data string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd.Add(data), msg.message}
}

func (msg Message) Marshal5424() []byte {
	var ret []byte
	prio := strconv.Itoa(int(msg.prio))
	ts := msg.timestamp.Format(RFC5424TimeStamp)
	sd := msg.sd.String()
	switch msg.message {
	case "":
		l := len(prio) + len(ts) + len(msg.hostname) + len(msg.appname) + len(msg.procid) + len(msg.msgid)
		l += len(sd)
		l += 9

		ret = make([]byte, 0, l)
		ret = append(ret, '<')
		ret = append(ret, []byte(prio)...)
		ret = append(ret, []byte{'>', '1', ' '}...)
		ret = append(ret, []byte(ts)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.hostname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.appname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.procid)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.msgid)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(sd)...)

	default:
		l := len(prio) + len(ts) + len(msg.hostname) + len(msg.appname) + len(msg.procid) + len(msg.msgid)
		l += len(sd) + len(msg.message)
		l += 10

		ret = make([]byte, 0, l)
		ret = append(ret, '<')
		ret = append(ret, []byte(prio)...)
		ret = append(ret, []byte{'>', '1', ' '}...)
		ret = append(ret, []byte(ts)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.hostname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.appname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.procid)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.msgid)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(sd)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.message)...)
	}
	return ret
}

func (msg Message) String() string {
	return string(msg.Marshal5424())
}

func CreateMessage(appname string, prio Priority, message string) Message {
	return EmptyMessage().AppName(appname).Priority(prio).LocalHost().Now().Msg(message)
}
