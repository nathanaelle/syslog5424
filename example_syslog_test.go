package syslog5424_test

import (
	"."
	"time"
	//"github.com/nathanaelle/syslog5424"
)

type someSD struct {
	Message string
	Errno   int
}

func ExampleSyslog() {
	syslog5424.Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	sl_conn := syslog5424.Dial("stdio", "stdout", syslog5424.T_LFENDED, -1)
	if sl_conn == nil {
		panic("no stderr available")
	}

	syslog, err := syslog5424.New(sl_conn, syslog5424.LOG_DAEMON|syslog5424.LOG_WARNING, "test app")
	if err != nil {
		panic(err.Error())
	}
	syslog.TestMode()

	conflog := syslog.SubSyslog("configuration")

	logger_info_conf := conflog.Channel(syslog5424.LOG_INFO).Logger("INFO : ")
	logger_err_conf := conflog.Channel(syslog5424.LOG_ERR).Logger("ERR : ")

	logger_info_conf.Print("doing some stuff")

	logger_err_conf.Print("doing some stuff")

	conflog.Channel(syslog5424.LOG_ERR).Log("another message", someSD{"some message", 42})
	time.Sleep(100 * time.Millisecond)
	sl_conn.End()
	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost test app/configuration 1234 - - ERR : doing some stuff
	//
	// <27>1 2014-12-20T14:04:00Z localhost test app/configuration 1234 - [someSD Message="some message" Errno="42"] another message

}
