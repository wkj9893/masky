package client

import (
	"fmt"
	"io"
	"net"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
)

func handleSocks(c *masky.Conn, config *Config) {
	defer c.Close()
	addr, err := socks.Handshake(c, config.Port)
	if err != nil {
		log.Error(err)
		return
	}
	var dst io.ReadWriteCloser
	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		log.Error(err)
		return
	}

	switch config.Mode {
	case DirectMode:
		if dst, err = masky.Dial(addr.String()); err != nil {
			log.Warn(err)
			return
		}
	case GlobalMode:
		if dst, err = masky.ConnectRemote(config.Addr); err != nil {
			log.Warn(err)
			return
		}
		if _, err := dst.Write(append([]byte{5}, addr...)); err != nil {
			log.Warn(err)
			return
		}
	case RuleMode:
		isocode, err := masky.Lookup(host, port)
		if err != nil {
			return
		}
		if isocode == "CN" {
			if dst, err = masky.Dial(net.JoinHostPort(host, port)); err != nil {
				//	we try to use proxy
				log.Info(fmt.Sprintf("fail to connect %v, use proxy instead", host))
				masky.Set(host, "")
				if dst, err = masky.ConnectRemote(config.Addr); err != nil {
					log.Warn(err)
					return
				}
				if _, err = dst.Write(append([]byte{5}, addr...)); err != nil {
					log.Warn(err)
					return
				}
			}
		} else {
			if dst, err = masky.ConnectRemote(config.Addr); err != nil {
				log.Warn(err)
				return
			}
			if _, err = dst.Write(append([]byte{5}, addr...)); err != nil {
				log.Warn(err)
				return
			}
		}
	default:
		panic("unknown mode")
	}
	defer dst.Close()
	masky.Relay(c, dst)
}
