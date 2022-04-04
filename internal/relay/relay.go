package relay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
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

	addr, id, err := auth(stream, *config)
	if err != nil {
		return err
	}
	c := masky.NewConn(stream)
	dst, err := masky.ConectRemote(addr, id)
	if err != nil {
		return err
	}

	head, err := c.Reader().Peek(1)
	if err != nil {
		return err
	}
	if head[0] == 5 { // socks
		masky.Relay(c, dst)
	} else { // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			return err
		}
		if req.Method == http.MethodConnect {
			if err = req.WriteProxy(dst); err != nil {
				return err
			}
			masky.Relay(c, dst)
		} else {
			if err = req.WriteProxy(dst); err != nil {
				return err
			}
			if _, err = io.Copy(c, dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func auth(stream quic.Stream, config Config) (string, uuid.UUID, error) {
	var id [16]byte
	_, err := stream.Read(id[:])
	if err != nil {
		return "", id, err
	}
	for _, v := range config.Proxies {
		if v.ID == id {
			return v.Server, v.ID, nil
		}
	}
	return "", id, errors.New("cannot authorzize user, fail to find uuid")
}
