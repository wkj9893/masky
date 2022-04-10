package masky

import (
	"crypto/tls"

	"github.com/google/uuid"
	"github.com/lucas-clemente/quic-go"
)

func ConectRemote(addr string, id uuid.UUID, tlsConf *tls.Config) (quic.Stream, error) {
	session, err := quic.DialAddrEarly(addr, tlsConf, nil)
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
