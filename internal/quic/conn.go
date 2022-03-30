package quic

import (
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/tls"
)

var (
	mu sync.Mutex

	addr string
	m    map[quic.EarlySession]quic.Stream
)

func SetAddr(s string) {
	mu.Lock()
	addr = s
	m = map[quic.EarlySession]quic.Stream{}
	mu.Unlock()
}

func ConectRemote() (quic.Stream, error) {
	if stream := find(); stream != nil {
		return stream, nil
	}
	session, err := quic.DialAddrEarly(addr, tls.ClientTLSConfig, masky.QuicConfig)
	if err != nil {
		return nil, err
	}
	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}
	mu.Lock()
	m[session] = stream
	mu.Unlock()
	return stream, nil
}

func find() quic.Stream {
	mu.Lock()
	defer mu.Unlock()
	for session, stream := range m {
		if !isActive(stream) {
			stream, err := session.OpenStream()
			if err != nil {
				delete(m, session)
				continue
			}
			m[session] = stream
			return stream
		}
	}
	return nil
}

func isActive(s quic.Stream) bool {
	select {
	case <-s.Context().Done():
		return false
	default:
		return true
	}
}
