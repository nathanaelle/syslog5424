package syslog5424 // import "github.com/nathanaelle/syslog5424"

import	(
	"os"
	"net"
)


type	(

	unixgram_receiver	struct {
		network		string
		address		string
		listener	*net.UnixConn
		pipeline	chan []byte
		end		chan struct{}
	}
)


func unixgram_coll(_, address string) Receiver {
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


func (r *unixgram_receiver) End() {
	close(r.end)
}


func (r *unixgram_receiver) SetTransport(_ Transport) {
}


func (r *unixgram_receiver) Receive() ([]byte, bool) {
	b, end := <- r.pipeline
	return b, end
}


func (r *unixgram_receiver) RunQueue(pipeline chan []byte) {
	defer	r.listener.Close()
	defer	close(pipeline)
	r.pipeline	= pipeline

	for {
		select {
		case <-r.end:
			return

		default:
			buffer := make([]byte, 65536)

			_,_,err := r.listener.ReadFrom(buffer)
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

			r.pipeline <- buffer[0:i]
		}
	}

}
