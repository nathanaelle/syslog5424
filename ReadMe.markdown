# Syslog5424

![License](http://img.shields.io/badge/license-Simplified_BSD-blue.svg?style=flat) [![Go Doc](http://img.shields.io/badge/godoc-syslog5424-blue.svg?style=flat)](http://godoc.org/github.com/nathanaelle/syslog5424) [![Build Status](https://travis-ci.org/nathanaelle/syslog5424.svg?branch=master)](https://travis-ci.org/nathanaelle/syslog5424)

## Example

```
import	(
	"github.com/nathanaelle/syslog5424"
)

type someSD struct{
	Message string
	Errno int
}

func main() {
	// create a connection to a server
	sl_conn,_:= syslog5424.Dial( "stdio", "stderr" )

	// create a syslog wrapper around the connection
	syslog,_ := syslog5424.New( sl_conn, syslog5424.LOG_DAEMON|syslog5424.LOG_WARNING, "test-app" )

	// create a channel for errors
	err_channel	:= syslog.Channel( syslog5424.LOG_ERR )

	// plug the golang log.Logger API to this channel
	logger_err := err_channel.Logger( "ERR : " )

	// log a message through the log.Logger
	logger_err.Print( "doing some stuff" )

	// log a message directly with some structured data
	err_channel.Log( "another message", someSD{ "some message", 42 } )
}

```
  * Example of client : [example_01-client_test.go](example_01-client_test.go)
  * Example of server :  [example_02-server_test.go](example_02-server_test.go)
  * Example of custom transport :  [example_03-custom_test.go](example_03-custom_test.go)

## Features

### Generic Features

  * [x] golang log.Logger compliant
  * [x] Handle multiple logging Channels
  * [x] Provide /dev/null Channel
  * [x] Extendable interfaces

### RFC 5424

  * [x] Encode RFC 5424 Message
  * [x] Decode RFC 5424 Message
  * [x] Encode Structured Data
  * [x] Decode Structured Data

### Networking / Communication

  * [x] Dial to a AF_UNIX datagram syslog server
  * [x] Dial to a AF_UNIX stream syslog server
  * [x] Dial to a TCP remote syslog server
  * [x] Accept to a AF_UNIX datagram syslog server
  * [x] Accept to a AF_UNIX stream syslog server
  * [x] Accept to a TCP remote syslog server

### Transport Encoding

  * [x] Unix Datagram Transport
  * [x] NULL terminated Transport
  * [x] LF terminated Transport
  * [x] RFC 5425 Transport

### Structured Data

  * [x] Encode Structured Data
  * [x] Decode Structured Data
  * [x] Encode Private Structured Data
  * [x] Decode Private Structured Data
  * [x] Decode Unknown Structured Data
  * [x] Structured Data Interface
  * [x] SDID Interface
  * [ ] SDIDLight Interface for Light Structured Data Support

### Structured Data Type

Source : [IANA syslog Structured Data ID Values](https://www.iana.org/assignments/syslog-parameters/syslog-parameters.xhtml#syslog-parameters-4)

  * [x] timeQuality (RFC 5424)
  * [ ] meta (RFC 5424)
  * [ ] origin (RFC 5424)
  * [ ] snmp (RFC 5675)
  * [ ] alarm (RFC 5674)
  * [ ] ssign (RFC 5848)
  * [ ] ssign-cert (RFC 5848)
  * [ ] PCNNode (RFC 6661)
  * [ ] PCNTerm (RFC 6661)

## License

2-Clause BSD

## Questions

### What is Syslog5424 ?

Syslog5424 is a library for coping with syslog messages through the log.Logger API.
Syslog5424 only produces syslog packets that are compatible with RFC 5424.
Those messages are not compatible with RFC 3164.

### What is Structured Data ?

The main point of the RFC 5424 is structured data.
This is a textual serialization of simple struct or map[string]string.
This serialization is _typed_ or _named_ and one text message can convey many Structured Data entries.
So This is a very pertinent way to mix *metrics*, *keywords* and human readable messages.

### What there is no support of UDP (RFC 5426) ?

System logging must be reliable for security audit of the log.
UDP is an unreliable protocol because UDP packet can be dropped, and neither the client nor the server will be informed of the missing data.

### Why remove parts of code about TLS ?

TLS is supported because the networing is implemented as interfaces.
but my idea of "security" is not compatible with maintaining duplicate code.

The requirement to support TLS are :

1. verify the certificate validity
2. verify the chain of trust to the root
3. verify OSCP staple if provided
4. check the OSCP's response from the CA
5. verify the CT with the OSCP's CT information and/or CT extra TLS header

so, you can :

1. Write your own code with the golang TLS stack (everything is provided through interfaces)
2. Wait for my implementation ( in [https://github.com/nathanaelle/pasnet](https://github.com/nathanaelle/pasnet) ) with the golang TLS stack wich will provide OCSP and Public Key verification

## Todo

  * Write documentation
  * Write comments
