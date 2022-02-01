package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
)

const (
	defaultStreamReceiveWindow     = 16777216 // 16 MB/s
	defaultConnectionReceiveWindow = 33554432 // 32 MB/s
	defaultMaxIncomingStreams      = 100      // maximum number of concurrent bidirectional streams
)

var DefaultQuicConfig = &quic.Config{
	HandshakeIdleTimeout:           time.Second,
	InitialStreamReceiveWindow:     defaultStreamReceiveWindow,
	MaxStreamReceiveWindow:         defaultStreamReceiveWindow,
	InitialConnectionReceiveWindow: defaultConnectionReceiveWindow,
	MaxConnectionReceiveWindow:     defaultConnectionReceiveWindow,
	MaxIncomingStreams:             defaultMaxIncomingStreams,
	KeepAlive:                      true,
}
