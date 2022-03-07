package masky

import (
	"time"

	"github.com/lucas-clemente/quic-go"
)

const (
	DefaultApplicationErrorCode = 0

	defaultMaxIncomingStreams = 1024 // maximum number of concurrent bidirectional streams
)

var (
	DefaultQuicConfig = &quic.Config{
		HandshakeIdleTimeout: time.Second,
		// https://www.rfc-editor.org/rfc/rfc9002.html#name-initial-and-minimum-congest
		InitialStreamReceiveWindow:     16 * (1 << 10), // 16 KB
		MaxStreamReceiveWindow:         10 * (1 << 20), // 10 MB
		InitialConnectionReceiveWindow: 16 * (1 << 10), // 16 KB
		MaxConnectionReceiveWindow:     25 * (1 << 20), // 25 MB
		MaxIncomingStreams:             defaultMaxIncomingStreams,
		EnableDatagrams:                true,
		KeepAlive:                      true,
	}
)
