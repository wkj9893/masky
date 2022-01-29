package http

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/geoip"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func HandleConn(c *masky.Conn, client *masky.Client) {
	local := true
	req, err := http.ReadRequest(c.Reader())
	if err != nil {
		return
	}

	var dst io.ReadWriteCloser
	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}

	switch client.Config().Mode {
	case masky.DirectMode:
		dst, err = masky.Dial(net.JoinHostPort(host, port))
		if err != nil {
			log.Warn(err)
			return
		}
	case masky.GlobalMode:
		dst, err = client.ConectRemote()
		if err != nil {
			log.Error(err)
			return
		}
		local = false
	case masky.RuleMode:
		ip, err := net.ResolveIPAddr("ip", req.URL.Hostname())
		if err != nil {
			log.Warn(err)
			return
		}
		isocode, err := geoip.Lookup(ip.IP)
		if err != nil {
			log.Warn(err)
			isocode = "CN"
		}
		if isocode == "CN" {
			if dst, err = masky.Dial(net.JoinHostPort(host, port)); err != nil {
				log.Warn(err)
				return
			}
		} else {
			if dst, err = client.ConectRemote(); err != nil {
				return
			}
			local = false
		}
	}

	if req.Method == http.MethodConnect {
		if !local {
			if err = req.WriteProxy(dst); err != nil {
				log.Error(err)
				return
			}
		}
		fmt.Fprintf(c, "%v %v \r\n\r\n", req.Proto, http.StatusOK)
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
