package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"./sdata"
	"bytes"
	"errors"
	"time"
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

var (
	ERR_Invalid     error = errors.New("Invalid")
	ERR_Pos0        error = errors.New("Pos 0 Found")
	ERR_PosNotFound error = errors.New("Pos Not Found")
	ERRImpossible   error = errors.New("NO ONE EXPECT THE RETURN OF SPANISH INQUISITION")
)

func search_next_sep(data, sep []byte) (pos int, err error) {
	pos = bytes.Index(data, sep)
	if pos > 0 {
		return
	}

	// pos at -1 or 0 means bad data
	err = ERR_PosNotFound
	if pos == 0 {
		err = ERR_Pos0
	}
	return
}

func Parse(data []byte) (msg MessageImmutable, err error) {
	sep_sp := []byte{' '}
	sep_brk := []byte{']'}
	err_msg := MessageImmutable{nil, nil, -1}

	msg = MessageImmutable{
		buffer: data,
		index:  make([]int, 0, 8),
		text:   -1,
	}

	parts := 0
	begin := 0
	for len(data) > 0 && parts < 6 {
		end, err := search_next_sep(data[begin:], sep_sp)
		if err != nil {
			//log.Printf("%s index %#d parts %d rest %s ", string(msg.buffer), msg.index, parts, data[begin:])
			return err_msg, err
		}
		msg.index = append(msg.index, begin+end)
		begin = begin + end + 1
		parts++
	}

	if len(data[begin:]) == 0 {
		return err_msg, ERR_Invalid
	}

	if len(data[begin:]) == 1 {
		if data[begin] != '-' {
			return err_msg, ERR_Invalid
		}
		msg.index = append(msg.index, begin+2)
		return
	}

	end, err := search_next_sep(data[begin:], sep_sp)
	if err != nil {
		return err_msg, err
	}

	if end == 1 {
		if data[begin] != '-' {
			return err_msg, ERR_Invalid
		}
		msg.index = append(msg.index, begin+1)
		msg.text = begin + 2
		return
	}

	t := begin
	for len(data[t:]) > 0 {
		end, err = search_next_sep(data[t:], sep_brk)

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
					return
				}
				continue
			}
			return err_msg, ERR_Invalid

		case ERR_PosNotFound:
			msg.text = begin + 1
			err = nil
			return

		default:
			return err_msg, err
		}
	}
	return err_msg, ERRImpossible
}

func (msg MessageImmutable) String() string {
	return string(msg.buffer)
}

func (msg MessageImmutable) Priority() (prio Priority) {
	begin := 0
	end := msg.index[0]
	(&prio).Unmarshal5424(msg.buffer[begin:end])
	return
}

func (msg MessageImmutable) TimeStamp() (ts time.Time) {
	begin := msg.index[0] + 1
	end := msg.index[1]
	ts, _ = time.Parse(RFC5424TimeStamp, string(msg.buffer[begin:end]))
	return

}

func (msg MessageImmutable) Hostname() (host string) {
	begin := msg.index[1] + 1
	end := msg.index[2]
	host = valid_host(string(msg.buffer[begin:end]))
	return

}

func (msg MessageImmutable) AppName() (app string) {
	begin := msg.index[2] + 1
	end := msg.index[3]
	app = valid_app(string(msg.buffer[begin:end]))
	return
}

func (msg MessageImmutable) ProcID() (procid string) {
	begin := msg.index[3] + 1
	end := msg.index[4]
	procid = valid_procid(string(msg.buffer[begin:end]))
	return
}

func (msg MessageImmutable) MsgID() (msgid string) {
	begin := msg.index[4] + 1
	end := msg.index[5]
	msgid = valid_msgid(string(msg.buffer[begin:end]))
	return
}

func (msg MessageImmutable) StructuredData() (lsd sdata.List) {
	lsd = sdata.List{}
	if len(msg.index[5:]) < 2 {
		return sdata.EmptyList()
	}

	if len(msg.index[5:]) == 2 && msg.buffer[msg.index[5]+1] == '-' {
		return sdata.EmptyList()
	}

	lsd_index := make([]int, len(msg.index[5:])-1)
	copy(lsd_index, msg.index[6:])
	begin := msg.index[5] + 1 // remember the first index have a space before like any previous field

	for len(lsd_index) > 0 {
		end := lsd_index[0]
		sd, ok := sdata.Parse(msg.buffer[begin:end])
		if ok {
			lsd = lsd.Add(sd)
		}
		begin = lsd_index[0]
		lsd_index = lsd_index[1:]
	}

	return
}

func (msg MessageImmutable) Message() (text string) {
	if msg.text < 1 {
		text = ""
		return
	}
	text = string(msg.buffer[msg.text:])
	return
}

func (m MessageImmutable) Writable() Message {
	return Message{m.Priority(), m.TimeStamp(), m.Hostname(), m.AppName(), m.ProcID(), m.MsgID(), m.StructuredData(), m.Message()}
}
