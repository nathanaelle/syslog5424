package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"log"
	"time"
)

type (
	Channel interface {
		io.Writer
		IsDevNull() bool
		AppName(string) Channel
		Msgid(string) Channel
		Logger(string) *log.Logger
		Log(string, ...interface{})
	}

	devnull struct {
	}

	msgChannel struct {
		devnull
		priority Priority
		hostname string
		pid      string
		appname  string
		msgid    string
		output   Conn
	}

	trueChannel struct {
		msgChannel
		priority Priority
		hostname string
		pid      string
		appname  string
		msgid    string
		output   Conn
	}
)

func (d *trueChannel) AppName(sup string) Channel {
	var appname string

	switch d.appname {
	case "-":
		appname = sup
	default:
		appname = d.appname + "/" + sup
	}

	return &trueChannel{
		priority: d.priority,
		hostname: d.hostname,
		pid:      d.pid,
		appname:  appname,
		msgid:    d.msgid,
		output:   d.output,
	}
}

func (d *trueChannel) Msgid(msgid string) Channel {
	return &msgChannel{
		priority: d.priority,
		hostname: d.hostname,
		pid:      d.pid,
		appname:  d.appname,
		msgid:    msgid,
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
	c.output.Send(Message{c.priority, time.Now(), c.hostname, c.appname, c.pid, c.msgid, emptyListSD, string(d)})
	return len(d), nil
}

func (c *msgChannel) Log(d string, sd ...interface{}) {
	switch len(sd) {
	case 0:
		c.output.Send(Message{c.priority, time.Now(), c.hostname, c.appname, c.pid, c.msgid, emptyListSD, string(d)})
	default:
		c.output.Send(Message{c.priority, time.Now(), c.hostname, c.appname, c.pid, c.msgid, listStructuredData(sd), string(d)})
	}
}

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
