package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"./sdata"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Message struct {
	prio      Priority
	timestamp time.Time
	hostname  string
	appname   string
	procid    string
	msgid     string
	sd        sdata.List
	message   string
}

// format of a RFC 5424 TimeStamp
const RFC5424TimeStamp string = "2006-01-02T15:04:05.999999Z07:00"

var hostname, _ = os.Hostname()

// Create a Message with the timestamp, hostname, appname, the priority and the message preset
func CreateMessage(appname string, prio Priority, message string) Message {
	return Message{prio, time.Now(), hostname, valid_app(appname), "-", "-", sdata.EmptyList(), strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Create a whole Message
func CreateWholeMessage(prio Priority, ts time.Time, host, app, pid, msgid, message string) Message {
	return Message{prio, ts, valid_host(host), valid_app(app), valid_procid(pid), valid_msgid(msgid), sdata.EmptyList(), strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Forge a whole Message
// hidden func because there is bypass
func forge_message(prio Priority, ts time.Time, host, app, pid, msgid, message string) Message {
	return Message{prio, ts, host, app, pid, msgid, sdata.EmptyList(), strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Create an empty Message
func EmptyMessage() Message {
	return Message{Priority(0), time.Unix(0, 0), "-", "-", "-", "-", sdata.EmptyList(), ""}
}

// Set the timestamp to time.Now()
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

// Set the timestamp from a time.Stamp string
func (msg Message) Timestamp(stamp string) Message {
	return Message{msg.prio, stamp_to_ts(stamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

func delta_boot_to_ts(boot_ts time.Time, s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return boot_ts.Add(time.Duration(nsec)*time.Nanosecond + time.Duration(sec)*time.Second)
}

// Set the timestamp from a time elapsed since boot time
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
func (msg Message) Host(host string) Message {
	return Message{msg.prio, msg.timestamp, valid_host(host), msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// set the app-name of a Message
func (msg Message) Time(ts time.Time) Message {
	return Message{msg.prio, ts, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// set the app-name of a Message
func (msg Message) AppName(appname string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, valid_app(appname), msg.procid, msg.msgid, msg.sd, msg.message}
}

// set the proc-id of a Message
func (msg Message) ProcID(procid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, valid_procid(procid), msg.msgid, msg.sd, msg.message}
}

// set the msg-id of a Message
func (msg Message) MsgID(msgid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, valid_msgid(msgid), msg.sd, msg.message}
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
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, strings.TrimRightFunc(message, unicode.IsSpace)}
}

//set the message part of a Message
func (msg Message) StructuredData(data ...sdata.StructuredData) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd.Add(data...), msg.message}
}

func (msg Message) Marshal5424() ([]byte, error) {
	var ret []byte
	prio := msg.prio.Marshal5424()
	ts := []byte(msg.timestamp.Format(RFC5424TimeStamp))
	sd, err := msg.sd.Marshal5424()
	if err != nil {
		return nil, err
	}
	switch msg.message {
	case "":
		l := len(prio) + len(ts) + len(msg.hostname) + len(msg.appname) + len(msg.procid) + len(msg.msgid)
		l += len(sd)
		l += 6

		ret = make([]byte, 0, l)
		ret = append(ret, prio...)
		ret = append(ret, ' ')
		ret = append(ret, ts...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.hostname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.appname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.procid)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.msgid)...)
		ret = append(ret, ' ')
		ret = append(ret, sd...)

	default:
		l := len(prio) + len(ts) + len(msg.hostname) + len(msg.appname) + len(msg.procid) + len(msg.msgid)
		l += len(sd) + len(msg.message)
		l += 7

		ret = make([]byte, 0, l)
		ret = append(ret, prio...)
		ret = append(ret, ' ')
		ret = append(ret, ts...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.hostname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.appname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.procid)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.msgid)...)
		ret = append(ret, ' ')
		ret = append(ret, sd...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.message)...)
	}
	return ret, nil
}

func (msg Message) String() (s string) {
	res, _ := msg.Marshal5424()
	s = string(res)
	return
}
