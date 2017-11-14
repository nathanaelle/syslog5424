package syslog5424 // import "github.com/nathanaelle/syslog5424"

// helper that provide the correct Listener when network is not known at compilation time
func GuessListener(network, address string) (Listener, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return TCPListener(network, address)
	case "unix":
		return UnixListener(address)
	case "unixgram":
		return UnixgramListener(address)

	}

	return nil, ErrorInvalidNetwork
}
