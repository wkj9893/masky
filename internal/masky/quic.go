package masky

import (
	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/tls"
)

func ConectRemote(addr string) (quic.Stream, error) {
	session, err := quic.DialAddrEarly(addr, tls.ClientTLSConfig, nil)
	if err != nil {
		return nil, err
	}
	return session.OpenStream()
}
