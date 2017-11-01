package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"log"
	"time"
)

var (
	// define the Now() function. only usefull in case of test or debug
	Now func() time.Time = time.Now
)

type (

	// Channel interface describe the common API of any logging channel.
	// a logging channel is a context set for any log message or structured data
	Channel interface {
		io.Writer
		IsDevNull() bool
		AppName(string) Channel
		Msgid(string) Channel
		Logger(string) *log.Logger
		Log(string, ...interface{})
	}

	// /dev/null Channel
	devnull struct {
	}

	// a msgChannel describe a channel with a MsgID already set
	msgChannel struct {
		devnull
		priority Priority
		hostname string
		pid      string
		appname  string
		msgid    string
		output   *Sender
	}

	//
	trueChannel struct {
		msgChannel
	}
)

func (d *trueChannel) AppName(sup string) Channel {
	var appname string

	switch string(d.appname) {
	case "-":
		appname = sup
	default:
		appname = d.appname + "/" + sup
	}

	return &trueChannel{msgChannel{
		priority: d.priority,
		hostname: d.hostname,
		pid:      d.pid,
		appname:  valid_app(appname),
		msgid:    d.msgid,
		output:   d.output,
	}}
}

func (d *trueChannel) Msgid(msgid string) Channel {
	return &msgChannel{
		priority: d.priority,
		hostname: d.hostname,
		pid:      d.pid,
		appname:  d.appname,
		msgid:    valid_msgid(msgid),
		output:   d.output,
	}
}

func (d *msgChannel) Logger(prefix string) *log.Logger {
	switch d.priority.Severity() {
	case LOG_DEBUG:
		return log.New(d, prefix, log.Lshortfile)
	default:
		return log.New(d, prefix, 0)
	}
}

func (d *msgChannel) IsDevNull() bool {
	return false
}

func (c *msgChannel) Write(d []byte) (int, error) {
	c.output.Send(forge_message(c.priority, Now(), c.hostname, c.appname, c.pid, c.msgid, string(d)))
	return len(d), nil
}

func (c *msgChannel) Log(d string, sd ...interface{}) {
	msg := forge_message(c.priority, Now(), c.hostname, c.appname, c.pid, c.msgid, string(d))

	if len(sd) > 0 {
		msg = msg.SetStructuredData(sd...)
	}

	c.output.Send(msg)
}

//	/dev/null Logger
// Nothing is logged
func (d *devnull) Logger(prefix string) *log.Logger {
	return log.New(d, prefix, 0)
}

func (dn *devnull) IsDevNull() bool {
	return true
}

func (dn *devnull) AppName(string) Channel {
	return dn
}

func (dn *devnull) Msgid(string) Channel {
	return dn
}

func (dn *devnull) Write(d []byte) (int, error) {
	return len(d), nil
}

func (dn *devnull) Log(_ string, _ ...interface{}) {
}
