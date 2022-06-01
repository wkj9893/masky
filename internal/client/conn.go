package client

import (
	"context"
	"io"
	"math/rand"

	"github.com/lucas-clemente/quic-go"
	"github.com/marten-seemann/webtransport-go"
	"golang.org/x/net/websocket"
)

func conectRemote(c *Config) (io.ReadWriteCloser, error) {
	api.RLock()
	defer api.RUnlock()
	i := api.index
	if i == 0 {
		i = rand.Intn(len(api.config.Proxies)-1) + 1
	}
	p := c.Proxies[i]
	switch p.Type {
	case "websocket":
		return websocket.Dial("wss://"+p.Server[0]+"/"+p.ID.String(), "", "http://localhost/")
	case "webtransport":
		var d webtransport.Dialer
		_, c, err := d.Dial(context.Background(), "https://"+p.Server[0]+"/"+p.ID.String(), nil)
		if err != nil {
			return nil, err
		}
		return c.OpenStream()
	default:
		c, err := quic.DialAddrEarly(p.Server[0], tlsConf, nil)
		if err != nil {
			return nil, err
		}
		conn, err := c.OpenStream()
		if err != nil {
			return nil, err
		}
		if _, err := conn.Write(p.ID[:]); err != nil {
			return nil, err
		}
		return conn, nil
	}
}
