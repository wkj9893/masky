package relay

import (
	"context"
	"io"

	"github.com/lucas-clemente/quic-go"
	"github.com/marten-seemann/webtransport-go"
	"golang.org/x/net/websocket"
)

func conectRemote(p *Proxy) (io.ReadWriteCloser, error) {
	switch p.Type {
	case "websocket":
		return websocket.Dial("wss://"+p.Server+"/"+p.ID.String(), "", "http://localhost/")
	case "webtransport":
		var d webtransport.Dialer
		_, c, err := d.Dial(context.Background(), "https://"+p.Server+"/"+p.ID.String(), nil)
		if err != nil {
			return nil, err
		}
		return c.OpenStream()
	default:
		c, err := quic.DialAddrEarly(p.Server, tlsConf, nil)
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
