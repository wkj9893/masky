package masky

import (
	"crypto/tls"
	"io"

	"github.com/lucas-clemente/quic-go"
)

func ConectRemote(addr string, tlsConf *tls.Config) (io.ReadWriteCloser, error) {
	session, err := quic.DialAddrEarly(addr, tlsConf, nil)
	if err != nil {
		return nil, err
	}
	return session.OpenStream()
	// if err != nil {
	// 	return nil, err
	// }
	// if _, err := stream.Write(id[:]); err != nil {
	// 	return nil, err
	// }
	// return stream, nil
}
