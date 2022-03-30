package client

import (
	"io"
	"net"

	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/quic"
	"github.com/wkj9893/masky/internal/socks"
)

func handleSocks(c *masky.Conn, mode Mode) error {
	defer c.Close()
	addr, err := socks.Handshake(c)
	if err != nil {
		return err
	}
	var dst io.ReadWriteCloser
	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return err
	}

	switch mode {
	case DirectMode:
		if dst, err = masky.Dial(addr.String()); err != nil {
			return err
		}
	case GlobalMode:
		if dst, err = quic.ConectRemote(); err != nil {
			return err
		}
		if _, err := dst.Write(append([]byte{5}, addr...)); err != nil {
			return err
		}
	case RuleMode:
		isocode, err := masky.Lookup(host, port)
		if err != nil {
			return nil
		}
		if isocode == "CN" {
			if dst, err = masky.Dial(net.JoinHostPort(host, port)); err != nil {
				return err
			}
		} else {
			if dst, err = quic.ConectRemote(); err != nil {
				return err
			}
			if _, err = dst.Write(append([]byte{5}, addr...)); err != nil {
				return err
			}
		}
	default:
		panic("unknown mode")
	}
	defer dst.Close()
	masky.Relay(c, dst)
	return nil
}
