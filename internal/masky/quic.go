package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
)

var DefaultQuicConfig = &quic.Config{
	HandshakeIdleTimeout: time.Second,
	KeepAlive:            true,
	// MaxIdleTimeout:       time.Hour,
}
