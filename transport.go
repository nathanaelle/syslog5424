package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"bufio"
	"bytes"
	"errors"
	"strconv"
)


type	(

	gen_t		struct{
		conn	Conn
	}

	// Encode frame in NULL terminated frame
	T_ZEROENDED	struct{
		gen_t
	}

	// Encode frame in LF terminated frame
	T_LFENDED	struct{
		gen_t
	}

	// Encode frame in RFC 5426 formated frame
	// RFC 5426 Format format is :
	// len([]byte) ' ' []byte
	T_RFC5426	struct{
		gen_t
	}

	Transport interface {
		Conn

		//
		SetConn(Conn)

		//
		Tokenize(io.ReadWriteCloser, chan<-[]byte)
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
func (t *gen_t) Read(d []byte) (int,error) {
	if t.conn == nil {
		return 0,errors.New("no Conn set")
	}

	return t.conn.Read(d)
}


func (t *gen_t) write_conn(data []byte) (int,error) {
	if t.conn == nil {
		return 0,errors.New("no Conn set")
	}

	p	:= 0
	t_len	:= len(data)
	for p < t_len {
		s,err := t.conn.Write(data[p:])
		p +=s
		if err == nil {
			continue
		}


		// TODO do some magic

	}

	return	t_len,nil
}


// split function for NULL terminated message
func (t *T_ZEROENDED) split(data []byte, atEOF bool) (int, []byte, error) {
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
func (t *T_ZEROENDED) Write(d []byte) (int,error) {
	return t.write_conn(append(d, byte(0) ))
}


func (t *T_ZEROENDED) Tokenize(conn io.ReadWriteCloser, pipeline chan<-[]byte) {
	scan	:= bufio.NewScanner(conn)
	scan.Split(t.split)

	for scan.Scan() {
		pipeline <- scan.Bytes()
	}

	conn.Close()
}



// split function for LF terminated message
func  (t *T_LFENDED) split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	// more data.
	return 0, nil, nil
}


// Write a LF terminated message
// see (Conn interface)[#Conn]
func (t *T_LFENDED) Write(d []byte) (int,error) {
	return t.write_conn(append(d, '\n' ))
}


func (t *T_LFENDED) Tokenize(conn io.ReadWriteCloser, pipeline chan<-[]byte) {
	scan	:= bufio.NewScanner(conn)
	scan.Split(t.split)

	for scan.Scan() {
		pipeline <- scan.Bytes()
	}

	conn.Close()
}


// split function for RFC 5426 message
func (t *T_RFC5426) split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	i := bytes.IndexByte(data, ' ')
	if i <= 0 {
		if len(data) < 10 {
			return 0, nil, nil
		}
		return 0, nil, errors.New("T_RFC5426 Split: no header len")
	}

	l, err := strconv.Atoi(string(data[0:i]))
	if err != nil {
		return 0, nil, errors.New("T_RFC5426 Split: invalid header len")
	}

	if len(data) < l {
		if atEOF {
			return 0, nil, errors.New("T_RFC5426 Split: incomplete message")
		}
		return 0, nil, nil
	}

	return i+l+1, data[i:i+l+1], nil
}


// Write a RFC 5426 formated message
// see (Conn interface)[#Conn]
func (t *T_RFC5426) Write(d []byte) (int,error) {
	l := len(d)
	h := []byte(strconv.Itoa(l))
	ret := make([]byte, l+len(h)+1)
	copy(ret[0:len(h)], h[:])
	ret[len(h)] = ' '
	copy(ret[len(h)+1:], d[:])

	return t.write_conn(ret)
}


func (t *T_RFC5426) Tokenize(conn io.ReadWriteCloser, pipeline chan<-[]byte) {
	scan	:= bufio.NewScanner(conn)
	scan.Split(t.split)

	for scan.Scan() {
		pipeline <- scan.Bytes()
	}

	conn.Close()
}
