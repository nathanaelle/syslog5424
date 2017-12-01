package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"errors"
)

type Priority int

const severityMask = 0x07
const facilityMask = 0xf8

// Severity.
const (
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

// Facility.
const (
	LOG_KERN Priority = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

var facility = map[string]Priority{
	"kern":     LOG_KERN,
	"user":     LOG_USER,
	"mail":     LOG_MAIL,
	"daemon":   LOG_DAEMON,
	"auth":     LOG_AUTH,
	"syslog":   LOG_SYSLOG,
	"lpr":      LOG_LPR,
	"news":     LOG_NEWS,
	"uucp":     LOG_UUCP,
	"cron":     LOG_CRON,
	"authpriv": LOG_AUTHPRIV,
	"ftp":      LOG_FTP,
	"local0":   LOG_LOCAL0,
	"local1":   LOG_LOCAL1,
	"local2":   LOG_LOCAL2,
	"local3":   LOG_LOCAL3,
	"local4":   LOG_LOCAL4,
	"local5":   LOG_LOCAL5,
	"local6":   LOG_LOCAL6,
	"local7":   LOG_LOCAL7,
}

var severity = map[string]Priority{
	"emerg":   LOG_EMERG,
	"alert":   LOG_ALERT,
	"crit":    LOG_CRIT,
	"err":     LOG_ERR,
	"warning": LOG_WARNING,
	"notice":  LOG_NOTICE,
	"info":    LOG_INFO,
	"debug":   LOG_DEBUG,
}

var severity_string = []string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"}
var facility_string = []string{"kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news",
	"uucp", "cron", "authpriv", "ftp", "-", "-", "-", "-",
	"local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7"}

func (p *Priority) Set(d string) error {
	pos := -1

	for i, c := range d {
		if c == '.' {
			pos = i
			break
		}
	}

	if pos == -1 {
		return errors.New("invalid syslog facility.severity [" + d + "]")
	}

	f, ok := facility[d[0:pos]]
	if !ok {
		return errors.New("invalid syslog facility.severity [" + d + "]")
	}

	s, ok := severity[d[pos+1:]]
	if !ok {
		return errors.New("invalid syslog facility.severity [" + d + "]")
	}

	*p = f | s

	return nil
}

func (p Priority) Facility() Priority {
	return p & facilityMask
}

func (p Priority) Severity() Priority {
	return p & severityMask
}

func (p Priority) String() string {
	return facility_string[p.Facility()>>3] + "." + severity_string[p.Severity()]
}

func (p Priority) Marshal5424() (data []byte, err error) {
	u := byte(int(p) % 10)
	d := byte(int(p)%100 - (int(p) % 10))
	c := byte(int(p) - (int(p) % 100))
	c = '0' + c/100
	d = '0' + d/10
	u = '0' + u/1

	if c > '0' {
		data = []byte{'<', c, d, u, '>', '1'}
		return
	}
	if d > '0' {
		data = []byte{'<', d, u, '>', '1'}
		return
	}
	data = []byte{'<', u, '>', '1'}

	return
}

func (p *Priority) Unmarshal5424(d []byte) error {
	s := len(d)
	if s < 4 {
		return errors.New("bad format : " + string(d))
	}
	if s > 6 {
		return errors.New("bad format : " + string(d))
	}
	if d[0] != '<' {
		return errors.New("bad format : " + string(d))
	}
	if d[s-1] != '1' {
		return errors.New("bad format : " + string(d))
	}
	if d[s-2] != '>' {
		return errors.New("bad format : " + string(d))
	}
	b := d[1 : s-2]
	t_p := int(0)
	for _, v := range b {
		if v < '0' || v > '9' {
			return errors.New("bad format : " + string(d))
		}
		t_p = t_p*10 + int(v-'0')
	}
	*p = Priority(t_p)

	return nil
}
