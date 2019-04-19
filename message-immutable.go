package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"bytes"
	"io"
	"time"

	"github.com/nathanaelle/syslog5424/v2/sdata"
	//	"log"
)

type (
	// MessageImmutable is an incoming message from a remote agent.
	// So it's a read only structure.
	MessageImmutable struct {
		buffer []byte
		index  []int
		text   int
	}
)

func searchNextSep(data, sep []byte) (pos int, err error) {
	pos = bytes.Index(data, sep)
	if pos > 0 {
		return
	}

	// pos at -1 or 0 means bad data
	err = ErrPosNotFound
	if pos == 0 {
		err = ErrPos0
	}
	return
}

func (msg MessageImmutable) String() string {
	return string(msg.buffer)
}

// Priority return the priority field of a MessageImmutable
func (msg MessageImmutable) Priority() (prio Priority) {
	begin := 0
	end := msg.index[0]
	(&prio).Unmarshal5424(msg.buffer[begin:end])
	return
}

// TimeStamp return the timestamp field of a MessageImmutable
func (msg MessageImmutable) TimeStamp() (ts time.Time) {
	begin := msg.index[0] + 1
	end := msg.index[1]
	ts, _ = time.Parse(RFC5424TimeStamp, string(msg.buffer[begin:end]))
	return

}

// Hostname return the hostname field of a MessageImmutable
func (msg MessageImmutable) Hostname() (host string) {
	begin := msg.index[1] + 1
	end := msg.index[2]
	host = validHost(string(msg.buffer[begin:end]))
	return

}

// AppName return the appname field of a MessageImmutable
func (msg MessageImmutable) AppName() (app string) {
	begin := msg.index[2] + 1
	end := msg.index[3]
	app = validApp(string(msg.buffer[begin:end]))
	return
}

// ProcID return the procid field of a MessageImmutable
func (msg MessageImmutable) ProcID() (procid string) {
	begin := msg.index[3] + 1
	end := msg.index[4]
	procid = validProcid(string(msg.buffer[begin:end]))
	return
}

// MsgID return the msgid field of a MessageImmutable
func (msg MessageImmutable) MsgID() (msgid string) {
	begin := msg.index[4] + 1
	end := msg.index[5]
	msgid = validMsgid(string(msg.buffer[begin:end]))
	return
}

// StructuredData return the Structured Data list of a MessageImmutable
func (msg MessageImmutable) StructuredData() (lsd sdata.List) {
	lsd = sdata.List{}
	if len(msg.index[5:]) < 2 {
		return sdata.EmptyList()
	}

	if len(msg.index[5:]) == 2 && msg.buffer[msg.index[5]+1] == '-' {
		return sdata.EmptyList()
	}

	lsdIndex := make([]int, len(msg.index[5:])-1)
	copy(lsdIndex, msg.index[6:])
	begin := msg.index[5] + 1 // remember the first index have a space before like any previous field

	for len(lsdIndex) > 0 {
		end := lsdIndex[0]
		sd, ok := sdata.Parse(msg.buffer[begin:end])
		if ok {
			lsd = lsd.Add(sd)
		}
		begin = lsdIndex[0]
		lsdIndex = lsdIndex[1:]
	}

	return
}

// Text return the text field of a MessageImmutable
func (msg MessageImmutable) Text() (text string) {
	if msg.text < 1 {
		text = ""
		return
	}
	text = string(msg.buffer[msg.text:])
	return
}

// Writable convert a MessageImmutable to Message
func (msg MessageImmutable) Writable() Message {
	return Message{msg.Priority(), msg.TimeStamp(), msg.Hostname(), msg.AppName(), msg.ProcID(), msg.MsgID(), msg.StructuredData(), msg.Text()}
}

// WriteTo implements io.WriterTo in MessageImmutable
func (msg MessageImmutable) WriteTo(w io.Writer) (n int64, err error) {
	in, err := w.Write(msg.buffer)
	n = int64(in)
	return
}

// Parse allow to parse a []byte and decode one ImmutableMessage
func Parse(data []byte, transport Transport, atEOF bool) (returnMsg MessageImmutable, rest []byte, mainErr error) {
	sepSp := []byte{' '}
	sepBrk := []byte{']'}

	if transport != nil {
		data, rest, mainErr = transport.PrefixStrip(data, atEOF)
		//log.Printf("P {%q} {%q} %v\n", data, rest, mainErr)
		if data == nil && rest == nil && mainErr == nil {
			return
		}
	}

	msg := MessageImmutable{
		buffer: data,
		index:  make([]int, 0, 8),
		text:   -1,
	}

	parts := 0
	begin := 0
	for len(data) > 0 && parts < 6 {
		end, err := searchNextSep(data[begin:], sepSp)
		if err != nil {
			//log.Printf("%s index %#d parts %d rest %s ", string(msg.buffer), msg.index, parts, data[begin:])
			mainErr = dispatchError(mainErr, err)
			return
		}
		msg.index = append(msg.index, begin+end)
		begin = begin + end + 1
		parts++
	}

	if len(data[begin:]) == 0 {
		mainErr = dispatchError(mainErr, ParseError{data, begin, "empty field expected"})
		return
	}

	if len(data[begin:]) == 1 {
		if data[begin] != '-' {
			mainErr = dispatchError(mainErr, ParseError{data, begin, "empty field expected"})
			return
		}
		msg.index = append(msg.index, begin+2)

		returnMsg, data, rest, mainErr = parseReturn(msg, transport, atEOF, data, rest)
		return
	}

	end, err := searchNextSep(data[begin:], sepSp)
	if err != nil {
		mainErr = dispatchError(mainErr, err)
		return
	}

	if end == 1 {
		if data[begin] != '-' {
			mainErr = dispatchError(mainErr, ParseError{data, begin, "empty structured data expected"})
			return
		}

		msg.index = append(msg.index, begin+1)
		msg.text = begin + 2
		returnMsg, data, rest, mainErr = parseReturn(msg, transport, atEOF, data, rest)
		return
	}

	t := begin
	for len(data[t:]) > 0 {
		end, err = searchNextSep(data[t:], sepBrk)

		switch err {
		case nil:
			if data[t+end-1] == '\\' {
				t = t + end + 1
				continue
			}

			if data[t+end-1] == '"' {
				msg.index = append(msg.index, t+end+1)
				begin = t + end + 1
				t = begin
				if len(data) <= begin {
					returnMsg, data, rest, mainErr = parseReturn(msg, transport, atEOF, data, rest)
					return
				}
				continue
			}

			mainErr = dispatchError(mainErr, ParseError{data, t + end - 1, `\\ or " expected`})
			return

		case ErrPosNotFound:
			msg.text = begin + 1

			returnMsg, data, rest, mainErr = parseReturn(msg, transport, atEOF, data, rest)
			return

		default:
			mainErr = dispatchError(mainErr, err)
			return
		}
	}
	mainErr = dispatchError(mainErr, ErrImpossible)
	return
}

func dispatchError(mainErr, err error) (returnErr error) {
	switch mainErr {
	case nil:
		returnErr = err
	default:
		returnErr = mainErr
	}
	return
}

func parseReturn(msg MessageImmutable, transport Transport, atEOF bool, oData, oRest []byte) (returnMsg MessageImmutable, data, rest []byte, mainErr error) {
	data = oData
	rest = oRest

	if transport == nil {
		returnMsg = msg
		return
	}

	if rest != nil {
		returnMsg = msg
		return
	}

	data, rest, mainErr = transport.SuffixStrip(data[msg.text:], atEOF)
	//log.Printf("RET\t{%q} {%q} %v\n", data, rest, mainErr)
	if data == nil && rest == nil && mainErr == nil {
		return
	}

	end := msg.text + len(data)
	if len(msg.buffer) > end {
		msg.buffer = msg.buffer[:end]

	}

	returnMsg = msg
	return
}
