package relay

import (
	"context"
	"errors"
	"fmt"

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
		go handleSession(s, config)
	}
}

func handleSession(s quic.EarlySession, config *Config) {
	stream, err := s.AcceptStream(context.Background())
	if err != nil {
		_ = s.CloseWithError(0, "")
		return
	}
	if err := handleStream(stream, config); err != nil {
		log.Error(err)
	}
}

func handleStream(stream quic.Stream, config *Config) error {
	defer stream.Close()
	var id [16]byte
	_, err := stream.Read(id[:])
	if err != nil {
		return err
	}
	for _, v := range config.Proxies {
		if v.ID == id {
			dst, err := masky.ConectRemote(v.Server, id)
			if err != nil {
				return err
			}
			go masky.Relay(stream, dst)
			return nil
		}
	}
	return errors.New("cannot authorzize user, fail to find uuid")
}
