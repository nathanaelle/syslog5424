package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"io"
	"log"
	"time"

	"github.com/nathanaelle/syslog5424/v2/sdata"
)

var (
	// Now permit to change the local default alias function for time.Now(). only usefull in case of test or debug
	Now = time.Now
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
		Log(string, ...sdata.StructuredData)
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
		appname:  validApp(appname),
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
		msgid:    validMsgid(msgid),
		output:   d.output,
	}
}

func (c *msgChannel) Logger(prefix string) *log.Logger {
	switch c.priority.Severity() {
	case LogDEBUG:
		return log.New(c, prefix, log.Lshortfile)
	default:
		return log.New(c, prefix, 0)
	}
}

func (c *msgChannel) IsDevNull() bool {
	return false
}

func (c *msgChannel) Write(d []byte) (int, error) {
	c.output.Send(forgeMessage(c.priority, Now(), c.hostname, c.appname, c.pid, c.msgid, string(d)))
	return len(d), nil
}

func (c *msgChannel) Log(d string, sd ...sdata.StructuredData) {
	msg := forgeMessage(c.priority, Now(), c.hostname, c.appname, c.pid, c.msgid, string(d))

	if len(sd) > 0 {
		msg = msg.StructuredData(sd...)
	}

	c.output.Send(msg)
}

//	/dev/null Logger
// Nothing is logged
func (dn *devnull) Logger(prefix string) *log.Logger {
	return log.New(dn, prefix, 0)
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

func (dn *devnull) Log(_ string, _ ...sdata.StructuredData) {
}
