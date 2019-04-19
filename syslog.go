package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"errors"
	"os"
	"strconv"
)

var (
	devNull = &devnull{}
)

// Syslog describes a sysloger
type Syslog struct {
	facility Priority
	hostname string
	pid      string
	appname  string
	channels []Channel
	output   *Sender
	minSev   int
}

// New create a new Syslog
func New(output *Sender, minPriority Priority, appname string) (syslog *Syslog, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}
	hostname = validHost(hostname)

	if appname == "" {
		err = errors.New("syslog.New needs a non empty appname")
		return
	}
	appname = validApp(appname)

	syslog = &Syslog{
		facility: minPriority.Facility(),
		hostname: hostname,
		pid:      strconv.Itoa(os.Getpid()),
		appname:  appname,
		output:   output,
		minSev:   int(minPriority.Severity()),
	}

	if syslog.pid == "" {
		syslog.pid = "-"
	}

	return syslog, nil
}

// TestMode prefill hostname and pid.
// this is only for test purpose
func (syslog *Syslog) TestMode() {
	syslog.hostname = "localhost"
	syslog.pid = "1234"
}

// Channel expose a custom channel for a priority
func (syslog *Syslog) Channel(sev Priority) Channel {
	if syslog.channels == nil {
		syslog.channels = []Channel{
			devNull, devNull, devNull, devNull,
			devNull, devNull, devNull, devNull,
		}

		for sev := 0; sev <= syslog.minSev; sev++ {
			syslog.channels[sev] = &trueChannel{msgChannel{
				priority: syslog.facility | Priority(sev),
				hostname: syslog.hostname,
				pid:      syslog.pid,
				appname:  syslog.appname,
				msgid:    "-",
				output:   syslog.output,
			}}
		}
	}

	return syslog.channels[sev.Severity()]
}

// SubSyslog create a syslog for a subAppName
// this allow to postfix to remplace the AppName with AppName/SubAppName
func (syslog *Syslog) SubSyslog(subAppName string) (sub *Syslog) {
	var appname string

	switch syslog.appname {
	case "-":
		appname = subAppName
	default:
		appname = syslog.appname + "/" + subAppName
	}

	sub = &Syslog{
		facility: syslog.facility,
		hostname: syslog.hostname,
		pid:      syslog.pid,
		appname:  appname,
		output:   syslog.output,
		channels: []Channel{
			devNull, devNull, devNull, devNull,
			devNull, devNull, devNull, devNull,
		},
		minSev: syslog.minSev,
	}

	for sev := 0; sev <= syslog.minSev; sev++ {
		sub.channels[sev] = &trueChannel{msgChannel{
			priority: syslog.facility | Priority(sev),
			hostname: syslog.hostname,
			pid:      syslog.pid,
			appname:  appname,
			msgid:    "-",
			output:   syslog.output,
		}}
	}

	return
}
