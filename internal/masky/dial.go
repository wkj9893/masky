package masky

import (
	"net"
	"time"
)

const (
	DefaultTcpTimeout = 1000 * time.Millisecond
)

func Dial(network string, address string) (net.Conn, error) {
	c, err := net.DialTimeout(network, address, DefaultTcpTimeout)
	if err != nil {
		return nil, err
	}
	return c, nil
}
