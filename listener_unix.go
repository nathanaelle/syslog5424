package syslog5424 // import "github.com/nathanaelle/syslog5424"

import	(
	"os"
	"net"
)


type	(
	unix_receiver	struct {
		network		string
		address		string
		listener	net.Listener
		end		chan struct{}
		transport	Transport
		pipeline	chan []byte
	}
)


func unix_coll(_, address string) Receiver {
	var err error

	r := new(unix_receiver)
	r.network	= "unix"
	r.address	= address
	r.end		= make(chan struct{})

	r.listener, err = net.ListenUnix("unix",  &net.UnixAddr { address, "unix" } )
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

		r.listener, err = net.ListenUnix("unix",  &net.UnixAddr { address, "unix" } )
	}

	return	r
}


func (r *unix_receiver) End() {
	close(r.end)
}


func (r *unix_receiver) SetTransport(t Transport) {
	r.transport	= t
}


func (r *unix_receiver) Receive() ([]byte, bool) {
	b, end := <- r.pipeline
	return b, end
}


func (r *unix_receiver) RunQueue(pipeline chan []byte) {
	defer	r.listener.Close()
	defer	close(pipeline)
	r.pipeline	= pipeline

	for {
		select {
		case <-r.end:
			return

		default:
			conn, err := r.listener.Accept()
			if err != nil {
				panic(err)
			}

			go r.transport.Tokenize( new_buffer(1<<18, buffer_read, conn), r.pipeline)
		}
	}

}
