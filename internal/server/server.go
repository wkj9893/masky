package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
	"github.com/wkj9893/masky/internal/tls"
	"gopkg.in/yaml.v3"
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
		go handleSession(s, config)
	}
}

func handleSession(s quic.EarlySession, config *Config) {
	stream, err := s.AcceptStream(context.Background())
	if err != nil {
		log.Warn("server error", err)
		return
	}
	if err := handleStream(stream, config); err != nil {
		log.Warn("server error", err)
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

func ParseConfig(name string) (*Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
