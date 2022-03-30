package client

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/quic"
)

func HandleHttp(c *masky.Conn, mode Mode) error {
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

	switch mode {
	case DirectMode:
		dst, err = masky.Dial(net.JoinHostPort(host, port))
		if err != nil {
			return err
		}
	case GlobalMode:
		dst, err = quic.ConectRemote()
		if err != nil {
			return err
		}
		local = false
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