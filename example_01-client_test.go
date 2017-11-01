package syslog5424

import (
	"log"
	"time"
)

type someSD struct {
	Message string
	Errno   int
}

func ExampleSyslogClient() {
	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	sl_conn, err := Dial("stdio", "stdout")
	if err != nil {
		log.Fatal(err)
	}

	syslog, err := New(sl_conn, LOG_DAEMON|LOG_WARNING, "test-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	conflog := syslog.SubSyslog("configuration")

	// using standard "log" API from golang
	logger_info_conf := conflog.Channel(LOG_INFO).Logger("INFO : ")
	logger_err_conf := conflog.Channel(LOG_ERR).Logger("ERR : ")

	// this is not logged because line 25 tell to syslog to log LOG_WARNING or higher
	logger_info_conf.Print("doing some stuff but not logged")

	logger_err_conf.Print("doing some stuff")

	// using internal API
	conflog.Channel(LOG_ERR).Log("another message with structured data", someSD{"some message", 42})

	// closing the connection and flushing all remaining logs
	sl_conn.End()

	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost test-app/configuration 1234 - - ERR : doing some stuff
	// <27>1 2014-12-20T14:04:00Z localhost test-app/configuration 1234 - [someSD Message="some message" Errno="42"] another message with structured data
}
