package http

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func HandleConn(c *masky.Conn, client *masky.Client) {
	defer c.Close()
	local := true
	req, err := http.ReadRequest(c.Reader())
	if err != nil {
		log.Error(err)
		return
	}

	var dst io.ReadWriteCloser
	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}

	switch client.Mode() {
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
		if masky.Lookup(host, port, client) == "CN" {
			if dst, err = masky.Dial(net.JoinHostPort(host, port)); err != nil {
				log.Warn(err)
				client.SetCache(host, "")
				return
			}
		} else {
			if dst, err = client.ConectRemote(); err != nil {
				log.Error(err)
				return
			}
			local = false
		}
	}

	defer dst.Close()
	if req.Method == http.MethodConnect {
		if !local {
			if err = req.WriteProxy(dst); err != nil {
				log.Error(err)
				return
			}
		}
		fmt.Fprintf(c, "%v %v \r\n\r\n", req.Proto, http.StatusOK)
		masky.Relay(c, dst)
		return
	}
	if err = req.WriteProxy(dst); err != nil {
		log.Error(err)
		return
	}
	if _, err = io.Copy(c, dst); err != nil {
		log.Error(err)
	}
}
