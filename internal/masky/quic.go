package masky

import (
	"github.com/google/uuid"
	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/tls"
)

func ConectRemote(addr string, id uuid.UUID) (quic.Stream, error) {
	session, err := quic.DialAddrEarly(addr, tls.ClientTLSConfig, nil)
	if err != nil {
		return nil, err
	}
	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}
	if _, err := stream.Write(id[:]); err != nil {
		return nil, err
	}
	return stream, nil
}
