package syslog5424 // import "github.com/nathanaelle/syslog5424"

import	(
	"os"
	"net"
	"time"
	"errors"
)


type	(

	unixgram_receiver	struct {
		network		string
		address		string
		listener	*net.UnixConn
		accepted	bool
		end		chan struct{}
	}


	fake_conn	struct {
		addr	*Addr
		buff	[]byte
		queue	chan []byte
		end	chan struct{}
	}

)


func unixgram_coll(_, address string) Listener {
	var err error

	r := new(unixgram_receiver)
	r.network	= "unixgram"
	r.address	= address
	r.end		= make(chan struct{})

	r.listener, err = net.ListenUnixgram("unixgram",  &net.UnixAddr { address, "unixgram" } )
	for err != nil {
		switch err.(type) {
			case *net.OpError:
				if err.(*net.OpError).Err.Error() != "bind: address already in use" {
					panic(err)
				}

			default:
				panic(err)
		}

		if _, r_err := os.Stat(address); r_err != nil {
			panic(err)
		}
		os.Remove(address)

		r.listener, err = net.ListenUnixgram("unixgram",  &net.UnixAddr { address, "unixgram" } )
	}
	return	r
}


func (r *unixgram_receiver) Close() (error) {
	close(r.end)
	return	r.listener.Close()
}


func (r *unixgram_receiver)Addr() net.Addr {
	return &Addr{ r.network, r.address }
}


// mimic an Accept
func (r *unixgram_receiver) Accept() (net.Conn, error) {
	if r.accepted {
		<-r.end
		return nil,errors.New("end")
	}

	r.accepted = true

	fc	:= &fake_conn{
		addr:	 &Addr{ r.network, r.address },
		queue:	make(chan []byte,1000),
		end:	make(chan struct{}),
	}

	go fc.run_queue(r.listener)

	return fc,nil
}


func (r *fake_conn)LocalAddr() net.Addr {
	return r.addr
}


func (r *fake_conn)RemoteAddr() net.Addr {
	return r.addr
}


func (r *fake_conn)SetDeadline(_ time.Time) error {
	return nil
}


func (r *fake_conn)SetReadDeadline(_ time.Time) error {
	return nil
}


func (r *fake_conn)SetWriteDeadline(_ time.Time) error {
	return nil
}


func (c *fake_conn) Redial() error {
	return nil
}


func (c *fake_conn) Flush() error {
	return nil
}


func (r *fake_conn) Close() error {
	close(r.end)
	return nil
}


func (r *fake_conn) Write(data []byte) (int, error) {
	return len(data),nil
}


func (r *fake_conn) Read(data []byte) (int, error) {
	if len(r.buff) == 0 {
		r.buff = <- r.queue
	}

	l_r	:= len(r.buff)
	l_d	:= len(data)
	if l_d <= l_r {
		copy(data[:], r.buff[0:l_d])
		r.buff = r.buff[l_d:]
		return l_d, nil
	}

	copy(data[0:l_r], r.buff[:])
	r.buff = nil

	return l_r, nil
}


func (r *fake_conn) run_queue(conn *net.UnixConn) {
	defer conn.Close()

	for {
		select {
		case <-r.end:
			return

		default:
			buffer := make([]byte, 65536)

			_,_,err := conn.ReadFrom(buffer)
			if err != nil {
				panic(err)
			}

			i := len(buffer)
			for i >0 {
				i--
				if buffer[i] != '0' {
					break
				}
			}

			r.queue <- buffer[0:i+2]
		}
	}

}
