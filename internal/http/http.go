package http

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/geoip"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/mode"
)

func HandleConn(c *masky.Conn, config masky.Config) {
	local := true
	req, err := http.ReadRequest(c.Reader())
	if err != nil {
		return
	}
	var dst io.ReadWriteCloser

	if config.Mode == mode.Direct {
		dst, err = masky.Dial("tcp", joinHostPort(req.URL.Hostname(), req.URL.Port()))
		if err != nil {
			return
		}
	} else if config.Mode == mode.Global {
		dst, err = masky.ConectRemote(config.Addr)
		if err != nil {
			return
		}
		local = false
	} else if config.Mode == mode.Rule {
		ip, err := net.LookupIP(req.URL.Hostname())
		if err != nil {
			return
		}
		isocode, err := geoip.Lookup(ip[0])
		if err != nil {
			log.Warn(err)
		}

		if isocode == "CN" {
			if dst, err = masky.Dial("tcp", joinHostPort(req.URL.Hostname(), req.URL.Port())); err != nil {
				return
			}
		} else {
			if dst, err = masky.ConectRemote(config.Addr); err != nil {
				return
			}
			local = false
		}
	} else {
		log.Error("Unknown Mode")
		return
	}

	if req.Method == http.MethodConnect {
		if local {
			fmt.Fprintf(c, "%v %v \r\n\r\n", req.Proto, http.StatusOK)
		} else {
			if err = req.WriteProxy(dst); err != nil {
				log.Error(err)
				return
			}
		}
		go masky.Copy(dst, c)
		go masky.Copy(c, dst)
		return
	}
	defer c.Close()
	if err = req.WriteProxy(dst); err != nil {
		log.Error(err)
		return
	}
	if _, err = io.Copy(c, dst); err != nil {
		log.Error(err)
	}
}

func joinHostPort(host string, port string) string {
	if port == "" {
		port = "80"
	}
	return net.JoinHostPort(host, port)
}
