package relay

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/tls"
)

func Run(config *Config) {
	tlsConf, err := tls.GenerateTLSConfig()
	if err != nil {
		panic(err)
	}
	l, err := quic.ListenAddrEarly(fmt.Sprintf(":%v", config.Port), tlsConf, nil)
	if err != nil {
		panic(err)
	}
	for {
		s, err := l.Accept(context.Background())
		if err != nil {
			panic(err)
		}
		go handleSession(s, config.Addrs)
	}
}

func handleSession(s quic.EarlySession, addr []string) {
	for {
		stream, err := s.AcceptStream(context.Background())
		if err != nil {
			_ = s.CloseWithError(0, "")
			return
		}
		if err := handleStream(stream, addr); err != nil {
			log.Error(err)
		}
	}
}

func handleStream(stream quic.Stream, addr []string) error {
	defer stream.Close()
	session, err := quic.DialAddr(addr[rand.Intn(len(addr))], tls.ClientTLSConfig, masky.QuicConfig)
	if err != nil {
		panic(err)
	}
	dst, err := session.OpenStream()
	if err != nil {
		panic(err)
	}
	go masky.Relay(stream, dst)
	return nil
}
