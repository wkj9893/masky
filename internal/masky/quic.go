package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
)

const (
	defaultStreamReceiveWindow     = 67108864 // 64 MB/s
	defaultConnectionReceiveWindow = 67108864 // 64 MB/s
	defaultMaxIncomingStreams      = 10000    // maximum number of concurrent bidirectional streams
)

var DefaultQuicConfig = &quic.Config{
	HandshakeIdleTimeout:           time.Second,
	MaxIdleTimeout:                 24 * time.Hour,
	InitialStreamReceiveWindow:     defaultStreamReceiveWindow,
	MaxStreamReceiveWindow:         defaultStreamReceiveWindow,
	InitialConnectionReceiveWindow: defaultConnectionReceiveWindow,
	MaxConnectionReceiveWindow:     defaultConnectionReceiveWindow,
	MaxIncomingStreams:             defaultMaxIncomingStreams,
	EnableDatagrams:                true,
	KeepAlive:                      true,
}
