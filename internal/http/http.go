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

	config := client.Config()
	switch config.Mode {
	case masky.DirectMode:
		dst, err = net.Dial("tcp", joinHostPort(req.URL.Hostname(), req.URL.Port()))
		if err != nil {
			return
		}
	case masky.GlobalMode:
		dst, err = client.ConectRemote()
		if err != nil {
			return
		}
		local = false
	case masky.RuleMode:
		ip, err := net.ResolveIPAddr("ip", req.URL.Hostname())
		if err != nil {
			return
		}
		isocode, err := geoip.Lookup(ip.IP)
		if err != nil {
			log.Warn(err)
			isocode = "CN"
		}
		if isocode == "CN" {
			if dst, err = net.Dial("tcp", joinHostPort(req.URL.Hostname(), req.URL.Port())); err != nil {
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
		go masky.Copy(c, dst)
		go masky.Copy(dst, c)
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
