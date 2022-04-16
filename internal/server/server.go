package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
	"github.com/wkj9893/masky/internal/tls"
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
	log.Info("start server successfully")
	for {
		s, err := l.Accept(context.Background())
		if err != nil {
			panic(err)
		}
		go handleConn(s)
	}
}

func handleConn(c quic.EarlyConnection) {
	stream, err := c.AcceptStream(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	if err := handleStream(stream, c.LocalAddr().String(), c.RemoteAddr().String()); err != nil {
		log.Error(err)
	}
}

func handleStream(stream quic.Stream, local string, remote string) error {
	defer stream.Close()
	c := masky.NewConn(stream)

	head, err := c.Reader().Peek(1)
	if err != nil {
		return err
	}
	switch head[0] {
	case 5: // socks
		if _, err := c.Reader().ReadByte(); err != nil {
			return err
		}
		addr, err := socks.ReadAddr(stream, make([]byte, 256))
		if err != nil {
			return err
		}
		log.Info(remote, "->", local, "->", addr)
		dst, err := masky.Dial(addr.String())
		if err != nil {
			return err
		}
		masky.Relay(stream, dst)
	default: // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			return err
		}
		log.Info(remote, "->", local, "->", req.Host)
		if req.Method == http.MethodConnect {
			dst, err := masky.Dial(req.Host)
			if err != nil {
				return err
			}
			masky.Relay(stream, dst)
		} else {
			req.RequestURI = ""
			resp, err := masky.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if err = resp.Write(stream); err != nil {
				return err
			}
		}
	}
	return nil
}
