package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
	"os"
	"strconv"
)

var (
	devNull *devnull = new(devnull)
)

type Syslog struct {
	facility Priority
	hostname string
	pid      string
	appname  string
	channels []Channel
	output   *Sender
	min_sev  int
}

func New(output *Sender, min_priority Priority, appname string) (syslog *Syslog, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}
	hostname = valid_host(hostname)

	if appname == "" {
		err = errors.New("syslog.New needs a non empty appname")
		return
	}
	appname = valid_app(appname)

	syslog = &Syslog{
		facility: min_priority.Facility(),
		hostname: hostname,
		pid:      strconv.Itoa(os.Getpid()),
		appname:  appname,
		output:   output,
		min_sev:  int(min_priority.Severity()),
	}

	if syslog.pid == "" {
		syslog.pid = "-"
	}


	return syslog, nil
}

func (syslog *Syslog) TestMode() {
	syslog.hostname = "localhost"
	syslog.pid = "1234"
}

func (syslog *Syslog) Channel(sev Priority) Channel {
	if syslog.channels == nil {
		syslog.channels = []Channel{
			devNull, devNull, devNull, devNull,
			devNull, devNull, devNull, devNull,
		}

		for sev := 0; sev <= syslog.min_sev; sev++ {
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

func (syslog *Syslog) SubSyslog(sub_appname string) (sub *Syslog) {
	var appname string

	switch syslog.appname {
	case "-":
		appname = sub_appname
	default:
		appname = syslog.appname + "/" + sub_appname
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
		min_sev: syslog.min_sev,
	}

	for sev := 0; sev <= syslog.min_sev; sev++ {
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
