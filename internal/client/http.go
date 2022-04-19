package client

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/wkj9893/masky/internal/masky"
)

func handleHttp(c *masky.Conn, config *Config) error {
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

	switch config.Mode {
	case DirectMode:
		dst, err = masky.Dial(net.JoinHostPort(host, port))
		if err != nil {
			return err
		}
	case GlobalMode:
		remote, id := getIndex()
		if dst, err = masky.ConectRemote(remote, tlsConf); err != nil {
			return err
		}
		if _, err := dst.Write(id[:]); err != nil {
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
			remote, id := getIndex()
			if dst, err = masky.ConectRemote(remote, tlsConf); err != nil {
				return err
			}
			if _, err := dst.Write(id[:]); err != nil {
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
