package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func handleHttp(c *masky.Conn, config *Config) {
	defer c.Close()
	req, err := http.ReadRequest(bufio.NewReader(c))
	if err != nil {
		log.Error(err)
		return
	}

	var dst io.ReadWriteCloser
	local := true
	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}

	switch config.Mode {
	case DirectMode:
		dst, err = masky.Dial(net.JoinHostPort(host, port))
		if err != nil {
			log.Warn(err)
			return
		}
	case GlobalMode:
		dst, err = masky.ConnectRemote(config.Addr)
		if err != nil {
			log.Warn(err)
			return
		}
		local = false
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
				local = false
			}
		} else {
			if dst, err = masky.ConnectRemote(config.Addr); err != nil {
				log.Warn(err)
				return
			}
			local = false
		}
	default:
		panic("unknown mode")
	}
	defer dst.Close()
	if req.Method == http.MethodConnect {
		if !local {
			if err = req.WriteProxy(dst); err != nil {
				log.Warn(err)
				return
			}
		}
		fmt.Fprintf(c, "%v %v \r\n\r\n", req.Proto, http.StatusOK)
	} else {
		if err = req.WriteProxy(dst); err != nil {
			log.Warn(err)
			return
		}
	}
	masky.Relay(c, dst)
}
