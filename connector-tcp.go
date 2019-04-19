package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"context"
	"net"
	"strconv"
	"time"
)

type (
	tcpConn struct {
		network, address string
	}
)

var resolver = &net.Resolver{
	PreferGo: true,
}

// TCPConnector is a dialer that forward to a local RFC5424 syslog receiver
func TCPConnector(network, address string) Connector {
	if network == "" || address == "" {
		return InvalidConnector{ErrEmptyNetworkAddress}
	}

	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return InvalidConnector{ErrInvalidNetwork}

	}

	return &tcpConn{network, address}
}

func (c *tcpConn) Connect() (conn WriteCloser, err error) {
	port := 514
	addr, portStr, err := net.SplitHostPort(c.address)
	if err != nil {
		return nil, err
	}

	if portStr != "" {
		port, err = strconv.Atoi(portStr)
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
		addr := &net.TCPAddr{IP: ip.IP, Port: port, Zone: ip.Zone}
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
