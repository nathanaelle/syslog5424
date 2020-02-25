/*

Package syslog5424 implements syslog RFC 5424 with a complete compatibility with the standard log.Logger API.


the simpliest way to use syslog5424 is :


	package main

	import	(
		"github.com/nathanaelle/syslog5424"
	)

	func main() {
		// create a connection the standard error
		sl_conn, _, _ := syslog5424.Dial( "stdio", "stderr:" )

		// create a syslog wrapper around the connection
		// the program is named "test-app" and you log to LogDAEMON facility at least at LogWARNING level
		syslog,_ := syslog5424.New( sl_conn, syslog5424.LogDAEMON|syslog5424.LogWARNING, "test-app" )

		// create a channel for the level LogERR
		err_channel	:= syslog.Channel( syslog5424.LogERR )

		// get a the *log.Logger for this channel with a prefix "ERR : "
		logger_err := err_channel.Logger( "ERR : " )

		// log a message through the log.Logger
		logger_err.Print( "doing some stuff" )
	}





*/
package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"
