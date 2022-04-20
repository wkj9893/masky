package relay

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/tls"
)

var (
	tlsConf = tls.ClientTLSConfig()
)

func Run(config *Config) {
	log.SetLogLevel(config.LogLevel)
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
		go handleConn(s, config)
	}
}

func handleConn(s quic.EarlyConnection, config *Config) {
	stream, err := s.AcceptStream(context.Background())
	if err != nil {
		log.Warn("relay error:", err)
		return
	}
	if err := handleStream(masky.NewConn(stream), config); err != nil {
		log.Warn("relay error:", err)
	}
}

func handleStream(c *masky.Conn, config *Config) error {
	defer c.Close()
	addr, err := auth(c, *config)
	if err != nil {
		return err
	}
	dst, err := masky.ConectRemote(addr, tlsConf)
	if err != nil {
		return err
	}
	masky.Relay(c, dst)
	return nil
}

func auth(c *masky.Conn, config Config) (string, error) {
	var i uuid.UUID
	id, err := c.Reader().Peek(16)
	if err != nil {
		return "", err
	}
	copy(i[:], id)
	for _, v := range config.Proxies {
		if v.ID == i {
			return v.Server, nil
		}
	}
	return "", errors.New("cannot authorzize user, fail to find uuid")
}
