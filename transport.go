package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bytes"
	"strconv"
)

type (
	// Encode frame in NULL terminated frame
	tZeroEnded struct{}

	// Encode frame in LF terminated frame
	tLFEnded struct{}

	// Encode frame in RFC 5425 formated frame
	// RFC 5425 Format format is :
	// len([]byte) ' ' []byte
	tRFC5425 struct{}

	tGuess struct{}

	Transport interface {
		// Set the sub conn where to write the transport-encoded data
		Encode([]byte) []byte

		// Decode the prefix in case of transport that use an encoding header
		PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error)

		// Decode the suffix in case of transport that use an encoding terminaison
		SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error)

		String() string
	}
)

var (
	// T_ZEROENDED is commonly used transport with "unix" and "unixgram"
	T_ZEROENDED Transport = tZeroEnded{}

	// T_LFENDED is commonly used transport with "tcp" "tcp4" and "tcp6"
	T_LFENDED Transport = tLFEnded{}

	// T_RFC5425 is performant transport specified in RFC 5425
	T_RFC5425 Transport = tRFC5425{}
)

func (t tZeroEnded) String() string {
	return "zero ended transport"
}

func (t tLFEnded) String() string {
	return "lf ended transport"
}

func (t tRFC5425) String() string {
	return "rfc 5425 transport"
}

func (t tZeroEnded) PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}

	return buffer, nil, nil
}

func (t tLFEnded) PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}
	return buffer, nil, nil
}

func (t tRFC5425) PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}

	if len(buffer) < 20 {
		return nil, nil, nil
	}

	sep_pos := bytes.IndexByte(buffer, ' ')
	if sep_pos <= 0 {
		return nil, nil, ERR_TRANSPORT_NOHEADER
	}

	lenMsg, err := strconv.Atoi(string(buffer[0:sep_pos]))
	if err != nil {
		return nil, nil, ERR_TRANSPORT_INVHEADER
	}

	start := sep_pos + 1
	lenBuff := start + lenMsg
	if len(buffer) < lenBuff {
		if atEOF {
			return buffer[start:], nil, ERR_TRANSPORT_INCOMPLETE
		}
		return nil, nil, nil
	}

	return buffer[start:lenBuff], buffer[lenBuff:], nil
}

func (t tZeroEnded) SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}

	if i := bytes.IndexByte(buffer, byte(0)); i >= 0 {
		return buffer[0:i], buffer[i+1:], nil
	}

	// at EOF act like \0 is implicit
	if atEOF {
		return buffer, nil, nil
	}

	return buffer, nil, nil
}

func (t tLFEnded) SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}

	if i := bytes.IndexByte(buffer, '\n'); i >= 0 {
		return buffer[0:i], buffer[i+1:], nil
	}

	if atEOF {
		return buffer, nil, ERR_TRANSPORT_INCOMPLETE
	}

	return buffer, nil, nil
}

func (t tRFC5425) SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}
	return buffer, nil, nil
}

// Write a NULL terminated message.
// see (Conn interface)[#Conn]
func (t tZeroEnded) Encode(d []byte) []byte {
	return append(d, byte(0))
}

// Write a LF terminated message
// see (Conn interface)[#Conn]
func (t tLFEnded) Encode(d []byte) []byte {
	return append(d, '\n')
}

// Write a RFC 5425 formated message
// see (Conn interface)[#Conn]
func (t tRFC5425) Encode(d []byte) []byte {
	l := len(d)
	h := []byte(strconv.Itoa(l))
	ret := make([]byte, l+len(h)+1)
	copy(ret[0:len(h)], h[:])
	ret[len(h)] = ' '
	copy(ret[len(h)+1:], d[:])

	return ret
}
