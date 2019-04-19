package syslog5424

import (
	"log"
	"time"
)

func ExampleDial_stdio() {
	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	slConn, _, err := Dial("stdio", "stdout:")
	if err != nil {
		log.Fatal(err)
	}

	syslog, err := New(slConn, LogDAEMON|LogWARNING, "test-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	conflog := syslog.SubSyslog("configuration")

	// using standard "log" API from golang
	loggerInfoConf := conflog.Channel(LogINFO).Logger("INFO : ")
	loggerErrConf := conflog.Channel(LogERR).Logger("ERR : ")

	// this is not logged because line 25 tell to syslog to log LogWARNING or higher
	loggerInfoConf.Print("doing some stuff but not logged")

	loggerErrConf.Print("doing some stuff")

	// using internal API
	conflog.Channel(LogERR).Log("another message with structured data", GenericSD(someSD{"some message", 42}))

	// closing the connection and flushing all remaining logs
	slConn.End()

	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost test-app/configuration 1234 - - ERR : doing some stuff
	// <27>1 2014-12-20T14:04:00Z localhost test-app/configuration 1234 - [someSD@32473 Message="some message" Errno="42"] another message with structured data
}

type someSD struct {
	Message string
	Errno   int
}

func (someSD) GetPEN() uint64 {
	return uint64(32473)
}

func (someSD) IsIANA() bool {
	return false
}

func (someSD) String() string {
	return "someSD@32473"
}
