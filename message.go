package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Message struct {
	priority  Priority
	timestamp time.Time
	hostname  string
	appName   string
	procID    string
	msgID     string
	sd        listStructuredData
	message   string
}

// format of a RFC 5424 TimeStamp
const RFC5424TimeStamp string = "2006-01-02T15:04:05.999999Z07:00"

var hostname, _ = os.Hostname()

func Parse(data []byte) (Message, error) {
	parts := bytes.SplitN(data, []byte{' '}, 8)

	switch len(parts) {
	case 7:
		prio := new(Priority)
		err := prio.Unmarshal5424(parts[0])
		if err != nil {
			return EmptyMessage(), errors.New("Wrong Priority : " + err.Error() + " : [" + string(parts[0]) + "]")
		}

		ts, err := time.Parse(RFC5424TimeStamp, string(parts[1]))
		if err != nil {
			return EmptyMessage(), errors.New("Wrong TS :" + string(parts[1]))
		}

		if string(parts[6]) == "-" {
			return Message{*prio, ts, string(parts[2]), string(parts[3]), string(parts[4]), string(parts[5]), emptyListSD, ""}, nil
		}

		return Message{*prio, ts, string(parts[2]), string(parts[3]), string(parts[4]), string(parts[5]), emptyListSD, ""}, nil

	case 8:
		prio := new(Priority)
		err := prio.Unmarshal5424(parts[0])
		if err != nil {
			return EmptyMessage(), errors.New("Wrong Priority : " + err.Error() + " : [" + string(parts[0]) + "]")
		}

		ts, err := time.Parse(RFC5424TimeStamp, string(parts[1]))
		if err != nil {
			return EmptyMessage(), errors.New("Wrong TS :" + string(parts[1]))
		}

		if string(parts[6]) == "-" {
			return Message{*prio, ts, string(parts[2]), string(parts[3]), string(parts[4]), string(parts[5]), emptyListSD, string(parts[7])}, nil
		}

		return Message{*prio, ts, string(parts[2]), string(parts[3]), string(parts[4]), string(parts[5]), emptyListSD, string(parts[7])}, nil

	default:
		return EmptyMessage(), errors.New("Wrong message :" + string(data))
	}
}

// Create a Message with the timestamp, hostname, appname, the priority and the message preset
func CreateMessage(appname string, prio Priority, message string) Message {
	return Message{prio, time.Now(), hostname, valid_app(appname), "-", "-", emptyListSD, strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Create a whole Message
func CreateWholeMessage(prio Priority, ts time.Time, host, app, pid, msgid, message string) Message {
	return Message{prio, ts, valid_host(host), valid_app(app), valid_procid(pid), valid_msgid(msgid), emptyListSD, strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Forge a whole Message
// hidden func because there is bypass
func forge_message(prio Priority, ts time.Time, host, app, pid, msgid, message string) Message {
	return Message{prio, ts, host, app, pid, msgid, emptyListSD, strings.TrimRightFunc(message, unicode.IsSpace)}
}

// Create an empty Message
func EmptyMessage() Message {
	return Message{Priority(0), time.Unix(0, 0), "-", "-", "-", "-", emptyListSD, ""}
}

// Set the timestamp to time.SetTimeNow()
func (msg Message) SetTimeNow() Message {
	msg.timestamp = time.Now()
	return msg
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

// Set the timestamp from a time.SetTimeStamp string
func (msg Message) SetTimeStamp(stamp string) Message {
	msg.timestamp = stamp_to_ts(stamp)
	return msg
}

func delta_boot_to_ts(boot_ts time.Time, s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return boot_ts.Add(time.Duration(nsec)*time.Nanosecond + time.Duration(sec)*time.Second)
}

// Set the timestamp from a time elapsed since boot time
func (msg Message) SetTimeDelta(boot_ts time.Time, s_sec string, s_nsec string) Message {
	msg.timestamp = delta_boot_to_ts(boot_ts, s_sec, s_nsec)
	return msg
}

func epoc_to_ts(s_sec string, s_nsec string) time.Time {
	sec, _ := strconv.ParseInt(s_sec, 10, 64)
	nsec, _ := strconv.ParseInt(s_nsec, 10, 64)

	return time.Unix(sec, nsec)
}

// Set the date of a Message with a epoch TimeStamp
func (msg Message) SetTimeEpoch(s_sec string, s_nsec string) Message {
	msg.timestamp = epoc_to_ts(s_sec, s_nsec)
	return msg
}

// Set the Hostname of a Message
func (msg Message) SetHostname(host string) Message {
	msg.hostname = valid_host(host)
	return msg
}

// Set the Timestamp of a Message
func (msg Message) SetTime(ts time.Time) Message {
	msg.timestamp = ts
	return msg
}

// Set the AppName of a Message
func (msg Message) SetAppName(appname string) Message {
	msg.appName = valid_app(appname)
	return msg
}

// Set the ProcID of a Message
func (msg Message) SetProcID(procid string) Message {
	msg.procID = valid_procid(procid)
	return msg
}

// Set the SetMsgID of a Message
func (msg Message) SetMsgID(msgid string) Message {
	msg.msgID = valid_msgid(msgid)
	return msg
}

// Set the Priority of a Message
func (msg Message) SetPriority(prio Priority) Message {
	msg.priority = prio
	return msg
}

// Set the Hostname as the value obtained with gethostbyname()
func (msg Message) SetLocalHost() Message {
	msg.hostname = hostname
	return msg
}

// Set the Message part of a Message
func (msg Message) SetMsg(message string) Message {
	msg.message = strings.TrimRightFunc(message, unicode.IsSpace)
	return msg
}

// Set the SetStructuredData part of a Message
func (msg Message) SetStructuredData(data ...interface{}) Message {
	msg.sd = msg.sd.Add(data...)
	return msg
}

func (msg Message) Priority() string {
	return msg.priority.String()
}

func (msg Message) Timestamp() time.Time {
	return msg.timestamp
}

func (msg Message) Hostname() string {
	return msg.hostname
}

func (msg Message) AppName() string {
	return msg.appName
}

func (msg Message) ProcID() string {
	return msg.procID
}

func (msg Message) MsgID() string {
	return msg.msgID
}

func (msg Message) StructuredDataString() string {
	return msg.sd.String()
}

func (msg Message) Message() string {
	return msg.message
}

func (msg Message) Marshal5424() []byte {
	var ret []byte
	prio := msg.priority.Marshal5424()
	ts := []byte(msg.timestamp.Format(RFC5424TimeStamp))
	sd := msg.sd.marshal5424()
	switch msg.message {
	case "":
		l := len(prio) + len(ts) + len(msg.hostname) + len(msg.appName) + len(msg.procID) + len(msg.msgID)
		l += len(sd)
		l += 6

		ret = make([]byte, 0, l)
		ret = append(ret, prio...)
		ret = append(ret, ' ')
		ret = append(ret, ts...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.hostname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.appName)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.procID)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.msgID)...)
		ret = append(ret, ' ')
		ret = append(ret, sd...)

	default:
		l := len(prio) + len(ts) + len(msg.hostname) + len(msg.appName) + len(msg.procID) + len(msg.msgID)
		l += len(sd) + len(msg.message)
		l += 7

		ret = make([]byte, 0, l)
		ret = append(ret, prio...)
		ret = append(ret, ' ')
		ret = append(ret, ts...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.hostname)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.appName)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.procID)...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.msgID)...)
		ret = append(ret, ' ')
		ret = append(ret, sd...)
		ret = append(ret, ' ')
		ret = append(ret, []byte(msg.message)...)
	}
	return ret
}

func (msg Message) String() string {
	return string(msg.Marshal5424())
}
