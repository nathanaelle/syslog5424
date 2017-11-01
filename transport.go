package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bytes"
	"errors"
	"strconv"
)

type (
	gen_t struct {
		conn Conn
	}

	// Encode frame in NULL terminated frame
	T_ZEROENDED struct {
		gen_t
	}

	// Encode frame in LF terminated frame
	T_LFENDED struct {
		gen_t
	}

	// Encode frame in RFC 5426 formated frame
	// RFC 5426 Format format is :
	// len([]byte) ' ' []byte
	T_RFC5426 struct {
		gen_t
	}

	Transport interface {
		Conn

		// Set the sub conn where to write the transport-encoded data
		SetConn(Conn)

		// see bufio.Scanner
		Split([]byte, bool) (int, []byte, error)

		String() string
	}
)

// see (Conn interface)[#Conn]
func (t *gen_t) Flush() error {
	if t.conn == nil {
		return nil
	}

	return t.conn.Flush()
}

func (t *gen_t) SetConn(c Conn) {
	t.conn = c
}

// see (Conn interface)[#Conn]
func (t *gen_t) Close() error {
	if t.conn == nil {
		return nil
	}

	return t.conn.Close()
}

// see (Conn interface)[#Conn]
func (t *gen_t) Redial() error {
	if t.conn == nil {
		return nil
	}

	return t.conn.Redial()
}

// see (Conn interface)[#Conn]
func (t *gen_t) Read(d []byte) (int, error) {
	if t.conn == nil {
		return 0, errors.New("no Conn set")
	}

	return t.conn.Read(d)
}

func (t *gen_t) write_conn(data []byte) (int, error) {
	if t.conn == nil {
		return 0, errors.New("no Conn set")
	}

	p := 0
	t_len := len(data)
	for p < t_len {
		s, err := t.conn.Write(data[p:])
		p += s
		if err == nil {
			continue
		}

		// TODO do some magic

	}

	return t_len, nil
}

func (t *gen_t) String() string {
	return "unknown transport"
}

func (t *T_ZEROENDED) String() string {
	return "zero ended transport"
}

func (t *T_LFENDED) String() string {
	return "lf ended transport"
}

func (t *T_RFC5426) String() string {
	return "rfc 5426 transport"
}

// split function for NULL terminated message
func (t *T_ZEROENDED) Split(data []byte, atEOF bool) (int, []byte, error) {
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

// Write a NULL terminated message.
// see (Conn interface)[#Conn]
func (t *T_ZEROENDED) Write(d []byte) (int, error) {
	return t.write_conn(append(d, byte(0)))
}

// split function for LF terminated message
func (t *T_LFENDED) Split(data []byte, atEOF bool) (int, []byte, error) {
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

// Write a LF terminated message
// see (Conn interface)[#Conn]
func (t *T_LFENDED) Write(d []byte) (int, error) {
	return t.write_conn(append(d, '\n'))
}

// split function for RFC 5426 message
func (t *T_RFC5426) Split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if len(data) < 20 {
		return 0, nil, nil
	}

	sep_pos := bytes.IndexByte(data, ' ')
	if sep_pos <= 0 {
		return 0, nil, errors.New("T_RFC5426 Split: no header len")
	}

	msg_len, err := strconv.Atoi(string(data[0:sep_pos]))
	if err != nil {
		return 0, nil, errors.New("T_RFC5426 Split: invalid header len")
	}

	start := sep_pos + 1
	buf_len := start + msg_len
	if len(data) < buf_len {
		if atEOF {
			return 0, nil, errors.New("T_RFC5426 Split: incomplete message")
		}
		return 0, nil, nil
	}

	return buf_len, data[start:buf_len], nil
}

// Write a RFC 5426 formated message
// see (Conn interface)[#Conn]
func (t *T_RFC5426) Write(d []byte) (int, error) {
	l := len(d)
	h := []byte(strconv.Itoa(l))
	ret := make([]byte, l+len(h)+1)
	copy(ret[0:len(h)], h[:])
	ret[len(h)] = ' '
	copy(ret[len(h)+1:], d[:])

	return t.write_conn(ret)
}
