package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
)

const DefaultApplicationErrorCode = 0

var QuicConfig = &quic.Config{
	HandshakeIdleTimeout: time.Second,
	EnableDatagrams:      true,
	KeepAlive:            true,
}
