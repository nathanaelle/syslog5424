# Syslog5424

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
  * [x] Handling multiple logging Channel
  * [x] Providing /dev/null Channel
  * [x] Extendable interfaces

### RFC 5424

  * [x] Encoding RFC 5424 Message
  * [x] Decoding RFC 5424 Message
  * [x] Encoding Structured Data
  * [x] Encoding Private Structured Data
  * [ ] Decoding Structured Data

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
  * [x] RFC 5426 Transport

## License

2-Clause BSD

## Questions

### What is Syslog5424 ?

Syslog5424 is a library for coping with syslog message through the log.Logger API.
Syslog5424 only produce syslog packet that are compatible with RFC 5424.
Those messages are not compatible with RFC 3164.

### What is Structured Data ?

The main point of the 5424 is structured data.
This is a textual serialization of simple struct or map[string]string.
This serialization is _typed_ or _named_ and one message can convey many Structured Data with one text message.
So This is a very pertinent way to mix *metrics* and *keywords* and human reading message.

### Why remove parts of code about TLS ?

TLS is supported because the networing is implemented as interfaces.
but my idea of "security" is not compatible with maintaining duplicate code.

so, you can :
  * 1. write your own code with the golang TLS stack (every things are provided through interfaces)
  * 2. wait my code ( in [https://github.com/nathanaelle/pasnet](https://github.com/nathanaelle/pasnet) ) with the golang TLS stack wich will provide OCSP and Public Key verification

## Todo

  * Write documentation [http://godoc.org/github.com/nathanaelle/syslog5424](http://godoc.org/github.com/nathanaelle/syslog5424)
  * Write comments
  * Clean ugly stuff
