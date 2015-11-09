package syslog5424 // import "github.com/nathanaelle/syslog5424"

import	(
	"net"
)

type	(
	tcp_receiver	struct {
		unix_receiver
	}
)


func tcp_coll(network, address string) Receiver {
	var err error

	r := new(tcp_receiver)
	r.network	= network
	r.address	= address
	r.end		= make(chan struct{})

	laddr, err	:= net.ResolveTCPAddr(network, address)
	if err != nil {
		panic(err)
	}

	r.listener, err = net.ListenTCP(network, laddr )
	if err != nil {
		panic(err)
	}

	return	r
}
