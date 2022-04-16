package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/tls"
)

func ConectRemote(addr string) (quic.Stream, error) {
	c, err := quic.DialAddrEarly(addr, tls.ClientTLSConfig, quicConfig)
	if err != nil {
		return nil, err
	}
	stream, err := c.OpenStream()
	if err != nil {
		return nil, err
	}
	return stream, nil
}

var quicConfig = &quic.Config{
	MaxIdleTimeout: time.Second,
}
