package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
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
	log.Info("start server successfully")
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

	if err := auth(stream, *config); err != nil {
		return err
	}
	c := masky.NewConn(stream)

	head, err := c.Reader().Peek(1)
	if err != nil {
		return err
	}
	if head[0] == 5 { // socks
		if _, err = c.Reader().ReadByte(); err != nil {
			return err
		}
		addr, err := socks.ReadAddr(c, make([]byte, 256))
		if err != nil {
			return err
		}
		dst, err := masky.Dial(addr.String())
		if err != nil {
			return err
		}
		masky.Relay(c, dst)
	} else { // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			return err
		}
		if req.Method == http.MethodConnect {
			dst, err := masky.Dial(req.Host)
			if err != nil {
				return err
			}
			masky.Relay(c, dst)
		} else {
			client := http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Timeout: 5 * time.Second,
			}
			req.RequestURI = ""
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if err = resp.Write(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func auth(stream quic.Stream, config Config) error {
	var id [16]byte
	_, err := stream.Read(id[:])
	if err != nil {
		return err
	}
	for _, v := range config.Proxies {
		if v.ID == id {
			return nil
		}
	}
	return errors.New("cannot authorzize user, fail to find uuid")
}
