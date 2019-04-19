package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"errors"
)

// Priority is the encoded form of the severity and facility of a Message
type Priority int

const severityMask = 0x07
const facilityMask = 0xf8

// Severity.
const (
	LogEMERG Priority = iota
	LogALERT
	LogCRIT
	LogERR
	LogWARNING
	LogNOTICE
	LogINFO
	LogDEBUG
)

// Facility.
const (
	LogKERN Priority = iota << 3
	LogUSER
	LogMAIL
	LogDAEMON
	LogAUTH
	LogSYSLOG
	LogLPR
	LogNEWS
	LogUUCP
	LogCRON
	LogAUTHPRIV
	LogFTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LogLOCAL0
	LogLOCAL1
	LogLOCAL2
	LogLOCAL3
	LogLOCAL4
	LogLOCAL5
	LogLOCAL6
	LogLOCAL7
)

var facility = map[string]Priority{
	"kern":     LogKERN,
	"user":     LogUSER,
	"mail":     LogMAIL,
	"daemon":   LogDAEMON,
	"auth":     LogAUTH,
	"syslog":   LogSYSLOG,
	"lpr":      LogLPR,
	"news":     LogNEWS,
	"uucp":     LogUUCP,
	"cron":     LogCRON,
	"authpriv": LogAUTHPRIV,
	"ftp":      LogFTP,
	"local0":   LogLOCAL0,
	"local1":   LogLOCAL1,
	"local2":   LogLOCAL2,
	"local3":   LogLOCAL3,
	"local4":   LogLOCAL4,
	"local5":   LogLOCAL5,
	"local6":   LogLOCAL6,
	"local7":   LogLOCAL7,
}

var severity = map[string]Priority{
	"emerg":   LogEMERG,
	"alert":   LogALERT,
	"crit":    LogCRIT,
	"err":     LogERR,
	"warning": LogWARNING,
	"notice":  LogNOTICE,
	"info":    LogINFO,
	"debug":   LogDEBUG,
}

var severityString = []string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"}
var facilityString = []string{"kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news",
	"uucp", "cron", "authpriv", "ftp", "-", "-", "-", "-",
	"local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7"}

// Set implement flag.Value interface
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

// Facility return the facility part of a priority
func (p Priority) Facility() Priority {
	return p & facilityMask
}

// Severity return the severity part of a priority
func (p Priority) Severity() Priority {
	return p & severityMask
}

func (p Priority) String() string {
	return facilityString[p.Facility()>>3] + "." + severityString[p.Severity()]
}

// Marshal5424 encode the priority to Syslog5424 format
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

// Unmarshal5424 decode the priority to Syslog5424 format
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
	tp := int(0)
	for _, v := range b {
		if v < '0' || v > '9' {
			return errors.New("bad format : " + string(d))
		}
		tp = tp*10 + int(v-'0')
	}
	*p = Priority(tp)

	return nil
}
