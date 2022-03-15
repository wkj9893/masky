package http

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func HandleConn(c *masky.Conn, client *masky.Client) error {
	defer c.Close()
	local := true
	req, err := http.ReadRequest(c.Reader())
	if err != nil {
		return err
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
			return err
		}
	case masky.GlobalMode:
		dst, err = client.ConectRemote()
		if err != nil {
			return err
		}
		local = false
	case masky.RuleMode:
		isocode, err := masky.Lookup(host, port, client)
		if err != nil {
			return nil
		}
		if isocode == "CN" {
			if dst, err = masky.Dial(net.JoinHostPort(host, port)); err != nil {
				log.Warn(fmt.Sprintf("fail to dial %v directly, set it using proxy instead", host))
				client.SetCache(host, "")
				if dst, err = client.ConectRemote(); err != nil {
					return err
				}
				local = false
			}
		} else {
			if dst, err = client.ConectRemote(); err != nil {
				return err
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
				return err
			}
		}
		fmt.Fprintf(c, "%v %v \r\n\r\n", req.Proto, http.StatusOK)
		masky.Relay(c, dst)
		return nil
	}
	if err = req.WriteProxy(dst); err != nil {
		return err
	}
	if _, err = io.Copy(c, dst); err != nil {
		return err
	}
	return nil
}
