package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/nathanaelle/syslog5424/v2/sdata"
)

// Message describe a mutable Syslog 5424 Message
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

// RFC5424TimeStamp is the format of a RFC 5424 TimeStamp
const RFC5424TimeStamp string = "2006-01-02T15:04:05.999999Z07:00"

var hostname, _ = os.Hostname()

// CreateMessage create a Message with the timestamp, hostname, appname, the priority and the message preset
func CreateMessage(appname string, prio Priority, message string) Message {
	return Message{prio, time.Now(), hostname, validApp(appname), "-", "-", sdata.EmptyList(), strings.TrimRightFunc(message, unicode.IsSpace)}
}

// CreateWholeMessage create a whole Message
func CreateWholeMessage(prio Priority, ts time.Time, host, app, pid, msgid, message string) Message {
	return Message{prio, ts, validHost(host), validApp(app), validProcid(pid), validMsgid(msgid), sdata.EmptyList(), strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Forge a whole Message
// hidden func because there is bypass
func forgeMessage(prio Priority, ts time.Time, host, app, pid, msgid, message string) Message {
	return Message{prio, ts, host, app, pid, msgid, sdata.EmptyList(), strings.TrimRightFunc(message, unicode.IsSpace)}
}

// EmptyMessage create an empty Message
func EmptyMessage() Message {
	return Message{Priority(0), time.Unix(0, 0), "-", "-", "-", "-", sdata.EmptyList(), ""}
}

// Now set the timestamp to time.Now()
func (msg Message) Now() Message {
	return Message{msg.prio, time.Now(), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

func stampToTs(stamp string) time.Time {
	now := time.Now()
	ts, _ := time.Parse(time.Stamp, stamp)
	year := now.Year()

	if now.Month() == 1 && ts.Month() == 12 {
		year--
	}

	return time.Date(year, ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), ts.Location())
}

// Timestamp set the timestamp from a time.Stamp string
func (msg Message) Timestamp(stamp string) Message {
	return Message{msg.prio, stampToTs(stamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

func deltaBootToTs(bootTs time.Time, strSec string, strNsec string) time.Time {
	sec, _ := strconv.ParseInt(strSec, 10, 64)
	nsec, _ := strconv.ParseInt(strNsec, 10, 64)

	return bootTs.Add(time.Duration(nsec)*time.Nanosecond + time.Duration(sec)*time.Second)
}

// Delta set the timestamp from a time elapsed since boot time
func (msg Message) Delta(bootTs time.Time, strSec string, strNsec string) Message {
	return Message{msg.prio, deltaBootToTs(bootTs, strSec, strNsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

func epochToTs(strSec string, strNsec string) time.Time {
	sec, _ := strconv.ParseInt(strSec, 10, 64)
	nsec, _ := strconv.ParseInt(strNsec, 10, 64)

	return time.Unix(sec, nsec)
}

// Epoch set the date of a Message with a epoch TimeStamp
func (msg Message) Epoch(strSec string, strNsec string) Message {
	return Message{msg.prio, epochToTs(strSec, strNsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// Host set the app-name of a Message
func (msg Message) Host(host string) Message {
	return Message{msg.prio, msg.timestamp, validHost(host), msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// Time set the timestamp of a Message
func (msg Message) Time(ts time.Time) Message {
	return Message{msg.prio, ts, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// AppName set the app-name of a Message
func (msg Message) AppName(appname string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, validApp(appname), msg.procid, msg.msgid, msg.sd, msg.message}
}

// ProcID set the proc-id of a Message
func (msg Message) ProcID(procid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, validProcid(procid), msg.msgid, msg.sd, msg.message}
}

// MsgID set the msg-id of a Message
func (msg Message) MsgID(msgid string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, validMsgid(msgid), msg.sd, msg.message}
}

// Priority set the priority of a Message
func (msg Message) Priority(prio Priority) Message {
	return Message{prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// LocalHost set the hostname as the value get with gethostbyname()
func (msg Message) LocalHost() Message {
	return Message{msg.prio, msg.timestamp, hostname, msg.appname, msg.procid, msg.msgid, msg.sd, msg.message}
}

// Msg set the message part of a Message
func (msg Message) Msg(message string) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd, strings.TrimRightFunc(message, unicode.IsSpace)}
}

// StructuredData set the message part of a Message
func (msg Message) StructuredData(data ...sdata.StructuredData) Message {
	return Message{msg.prio, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.sd.Add(data...), msg.message}
}

// Marshal5424 encode a message to the syslog 5424 format
func (msg Message) Marshal5424() ([]byte, error) {
	var ret []byte
	prio, err := msg.prio.Marshal5424()
	if err != nil {
		return nil, err
	}

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
