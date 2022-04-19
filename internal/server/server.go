package server

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
		go handleConn(s, config)
	}
}
func handleConn(c quic.EarlyConnection, config *Config) {
	stream, err := c.AcceptStream(context.Background())
	if err != nil {
		log.Warn("server error:", err)
		return
	}
	if err := handleStream(masky.NewConn(stream), c.LocalAddr().String(), c.RemoteAddr().String(), config); err != nil {
		log.Warn("server error:", err)
	}
}

func handleStream(c *masky.Conn, local string, remote string, config *Config) error {
	defer c.Close()
	if err := auth(c, *config); err != nil {
		return err
	}
	head, err := c.Reader().Peek(1)
	if err != nil {
		return err
	}
	fmt.Println(head[0])
	switch head[0] {
	case 5: // socks
		if _, err := c.Reader().ReadByte(); err != nil {
			return err
		}
		addr, err := socks.ReadAddr(c, make([]byte, 256))
		if err != nil {
			return err
		}
		log.Info(remote, "->", local, "->", addr)
		dst, err := masky.Dial(addr.String())
		if err != nil {
			return err
		}
		masky.Relay(c, dst)
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
			masky.Relay(c, dst)
		} else {
			req.RequestURI = ""
			resp, err := masky.DefaultClient.Do(req)
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

func auth(c *masky.Conn, config Config) error {
	var id uuid.UUID
	_, err := io.ReadFull(c, id[:])
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
