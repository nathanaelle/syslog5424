package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bytes"
	"errors"
	"strconv"
)

type (
	// Encode frame in NULL terminated frame
	t_ZEROENDED struct{}

	// Encode frame in LF terminated frame
	t_LFENDED struct{}

	// Encode frame in RFC 5425 formated frame
	// RFC 5425 Format format is :
	// len([]byte) ' ' []byte
	t_RFC5425 struct{}

	T_GUESS struct{}

	Transport interface {
		// Set the sub conn where to write the transport-encoded data
		Encode([]byte) []byte

		// see bufio.Scanner
		//		Split([]byte, bool) (int, []byte, error)
		PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error)
		SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error)

		String() string
	}
)

var (
	T_ZEROENDED Transport = t_ZEROENDED{}
	T_LFENDED   Transport = t_LFENDED{}
	T_RFC5425   Transport = t_RFC5425{}
)

func (t t_ZEROENDED) String() string {
	return "zero ended transport"
}

func (t t_LFENDED) String() string {
	return "lf ended transport"
}

func (t t_RFC5425) String() string {
	return "rfc 5425 transport"
}

func (t t_ZEROENDED) PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}

	return buffer, nil, nil
}

func (t t_LFENDED) PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}
	return buffer, nil, nil
}

var (
	ERR_TRANSPORT_NOHEADER  error = errors.New("T_RFC5425 Split: no header len")
	ERR_TRANSPORT_INVHEADER error = errors.New("T_RFC5425 Split: invalid header len")
)

func (t t_RFC5425) PrefixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
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

	msg_len, err := strconv.Atoi(string(buffer[0:sep_pos]))
	if err != nil {
		return nil, nil, ERR_TRANSPORT_INVHEADER
	}

	start := sep_pos + 1
	buf_len := start + msg_len
	if len(buffer) < buf_len {
		if atEOF {
			return buffer[start:], nil, ERR_TRANSPORT_INCOMPLETE
		}
		return nil, nil, nil
	}

	return buffer[start:buf_len], buffer[buf_len:], nil
}

func (t t_ZEROENDED) SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
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

func (t t_LFENDED) SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
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

func (t t_RFC5425) SuffixStrip(buffer []byte, atEOF bool) (data, rest []byte, err error) {
	if buffer == nil || len(buffer) == 0 {
		return nil, nil, nil
	}
	return buffer, nil, nil
}

// Write a NULL terminated message.
// see (Conn interface)[#Conn]
func (t t_ZEROENDED) Encode(d []byte) []byte {
	return append(d, byte(0))
}

// Write a LF terminated message
// see (Conn interface)[#Conn]
func (t t_LFENDED) Encode(d []byte) []byte {
	return append(d, '\n')
}

/*
// split function for NULL terminated message
func (t t_ZEROENDED) Split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, byte(0)); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// TODO need to detect the non zero ended message here

	if atEOF {
		return len(data), data, nil
	}

	// more data.
	return 0, nil, nil
}

// split function for LF terminated message
func (t t_LFENDED) Split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		//return len(data), data, nil
		return 0, nil, errors.New("T_LFENDED Split: incomplete message")
	}

	// more data.
	return 0, nil, nil
}

// split function for RFC 5425 message
func (t t_RFC5425) Split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if len(data) < 20 {
		return 0, nil, nil
	}

	sep_pos := bytes.IndexByte(data, ' ')
	if sep_pos <= 0 {
		return 0, nil, errors.New("T_RFC5425 Split: no header len")
	}

	msg_len, err := strconv.Atoi(string(data[0:sep_pos]))
	if err != nil {
		return 0, nil, errors.New("T_RFC5425 Split: invalid header len")
	}

	start := sep_pos + 1
	buf_len := start + msg_len
	if len(data) < buf_len {
		if atEOF {
			return 0, nil, errors.New("T_RFC5425 Split: incomplete message")
		}
		return 0, nil, nil
	}

	return buf_len, data[start:buf_len], nil
}

*/

// Write a RFC 5425 formated message
// see (Conn interface)[#Conn]
func (t t_RFC5425) Encode(d []byte) []byte {
	l := len(d)
	h := []byte(strconv.Itoa(l))
	ret := make([]byte, l+len(h)+1)
	copy(ret[0:len(h)], h[:])
	ret[len(h)] = ' '
	copy(ret[len(h)+1:], d[:])

	return ret
}
