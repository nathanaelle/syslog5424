package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"io"
	"sync"
	"errors"
)


type (
	t_buffer	int

	buffer	struct {
		l	*sync.Mutex
		conn	io.ReadWriteCloser
		buff	[]byte
		size	int
		r	int
		w	int
		t	t_buffer
	}
)


const	(
	buffer_unknown	t_buffer	= iota
	buffer_read
	buffer_write
)


func new_buffer(l int, t t_buffer, c io.ReadWriteCloser) *buffer {
	return &buffer {
		l:	new(sync.Mutex),
		conn:	c,
		buff:	make([]byte,l),
		size:	l,
		t:	t,
	}
}


func (b *buffer) SetConn(c io.ReadWriteCloser) {
	b.conn = c
}



// io.Closer
func (b *buffer) Close() error {
	return b.Flush()

/*	if err := b.Flush(); err != nil {
		return err
	}

	return b.conn.Close() */
}


// io.Reader
func (b *buffer) Read(data []byte) (int, error) {
	if b.t != buffer_read {
		return 0, errors.New("not read buffer")
	}

	b.l.Lock()
	defer b.l.Unlock()

	l := len(data)
	d	:= b.r - b.w

	// if enough data in buffer, read it
	if (b.w + l) < b.r {
		copy(data[:], b.buff[b.w:b.w+l])
		b.w += l
		return l, nil
	}

	if d > 0 {
		copy(data[0:d], b.buff[b.w:b.r])
	}
	b.w	= 0
	b.r	= 0

	n, err := b.conn.Read(b.buff[:])

	// clearly not enough data
	if n == 0 {
		b.w = 0
		b.r = 0
		return d, err
	}

	// not enough data
	if d+n <= l {
		copy(data[d:], b.buff[0:n])
		b.w = 0
		b.r = 0
		return d+n, err
	}

	b.r = n
	b.w = l-d
	copy(data[d:], b.buff[0:b.w])

	// there were enough data, we will care of the error at the end of the buffer
	return l, nil
}


// io.Writer
func (b *buffer) Write(data []byte) (int, error) {
	if b.t != buffer_write {
		return 0, errors.New("not write buffer")
	}

	b.l.Lock()
	defer b.l.Unlock()

	l := len(data)
	// if not enough empty space in buffer, flush it
	if (b.w + l) > b.size {
		if err := b.true_flush(); err != nil {
			return 0,err
		}
	}

	copy(b.buff[b.w:],data[:])
	b.w+=l

	return l, nil
}


// flush the pending write
func (b *buffer) Flush() (err error) {
	if b.t != buffer_write {
		return errors.New("not write buffer")
	}

	b.l.Lock()
	defer b.l.Unlock()

	return b.true_flush()
}


func (b *buffer) true_flush() (err error) {
	t_n	:= 0
	for b.r < b.w {
		t_n, err = b.conn.Write(b.buff[b.r:b.w])
		b.r += t_n
		if err != nil {
			return
		}
	}

	b.w = 0
	b.r = 0

	return nil
}
