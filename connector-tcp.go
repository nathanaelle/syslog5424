package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"context"
	"net"
	"strconv"
	"time"
)

type (
	tcp_conn struct {
		network, address string
	}
)

var resolver = &net.Resolver{
	PreferGo: true,
}

// dialer that forward to a local RFC5424 syslog receiver
func TCPConnector(network, address string) Connector {
	if network == "" || address == "" {
		return InvalidConnector{ErrorEmptyNetworkAddress}
	}

	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return InvalidConnector{ErrorInvalidNetwork}

	}

	return &tcp_conn{network, address}
}

func (c *tcp_conn) Connect() (conn WriteCloser, err error) {
	port := 514
	addr, s_port, err := net.SplitHostPort(c.address)
	if err != nil {
		return nil, err
	}

	if s_port != "" {
		port, err = strconv.Atoi(s_port)
		if err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rawips, err := resolver.LookupIPAddr(ctx, addr)
	cancel()
	if err != nil {
		return nil, err
	}

	ips := make([]net.IPAddr, 0, len(rawips))
	switch c.network {
	case "tcp4":
		for _, ip := range rawips {
			if ip.IP.To4() != nil {
				ips = append(ips, ip)
			}
		}
	case "tcp6":
		for _, ip := range rawips {
			if ip.IP.To4() == nil {
				ips = append(ips, ip)
			}
		}
	default:
		ips = rawips
	}

	var contcp *net.TCPConn
	for _, ip := range ips {
		addr := &net.TCPAddr{ip.IP, port, ip.Zone}
		contcp, err = net.DialTCP(c.network, nil, addr)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	contcp.CloseRead()
	contcp.SetNoDelay(false)
	contcp.SetKeepAlive(true)
	contcp.SetKeepAlivePeriod(2 * time.Minute)
	contcp.SetLinger(-1)

	conn = contcp
	return
}
