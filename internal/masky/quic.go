package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
)

var QuicConfig = &quic.Config{
	HandshakeIdleTimeout: time.Second,
	MaxIdleTimeout:       5 * time.Minute,
	EnableDatagrams:      true,
}
