# Syslog5424

## Example

```
import	(
	"os"
	"github.com/nathanaelle/syslog5424"
)

func main() {

	syslog := syslog5424.New( os.Stderr, syslog5424.LOG_DAEMON|syslog5424.LOG_WARNING, "test app" )

	conflog := syslog.SubSyslog("configuration")

	logger_info_conf := conflog.Channel(syslog5424.LOG_INFO).Logger("")
	logger_error_conf := conflog.Channel(syslog5424.LOG_ERR).Logger("ERROR :")

	logger_info_conf.Print("doing some stuff")
	logger_error_conf.Printf("%#v", struct{ message string, errno int }{ "some message", 42 })

	conflog.Channel(syslog5424.LOG_INFO).Log("another message", )
}

```

## What is Syslog5424 ?

Syslog5424 is a library for coping with syslog message through the log.Logger API.
Syslog5424 only produce syslog packet that are compatible with RFC 5424.
Those messages are not compatible with RFC 3164.

### Structured Data

The main point of the 5424 is structured data.
This is a textual serialization of simple struct or map[string]string.
This serialization is _typed_ or _named_ and one message can convey many Structured Data with one text message.
So This is a very pertinent way to mix *metrics* and *keywords* and human reading message.

## Features

  * [x] Encoding Structured Data
  * [x] Encoding RFC 5424 Message
  * [x] Encoding Private Structured Data
  * [ ] Decoding Structured Data
  * [ ] Decoding RFC 5424 Message
  * [x] Handling multi channels
  * [x] Dial to a local unixdgram syslog server
  * [ ] Dial to a TCP remote syslog server
  * [ ] Dial to a TLS remote syslog server
  * [x] Unix Datagram Transport
  * [ ] LF separated transport
  * [ ] RFC 5426 Transport

## License
2-Clause BSD

## Todo

  * Write documentation [http://godoc.org/github.com/nathanaelle/syslog5424](http://godoc.org/github.com/nathanaelle/syslog5424)
  * Write comments
  * Clean some ugly stuff
