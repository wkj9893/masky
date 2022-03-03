package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
	"github.com/wkj9893/masky/internal/tls"
)

var config masky.ServerConfig

func init() {
	// default config
	config = masky.ServerConfig{
		Port:     2022,
		Password: "masky",
		LogLevel: log.InfoLevel,
	}
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
}

func main() {
	tlsConf, err := tls.GenerateTLSConfig()
	if err != nil {
		panic(err)
	}
	l, err := quic.ListenAddr(fmt.Sprintf(":%v", config.Port), tlsConf, masky.DefaultQuicConfig)
	if err != nil {
		panic(err)
	}
	for {
		s, err := l.Accept(context.Background())
		if err != nil {
			panic(err)
		}
		if err := auth(s); err != nil {
			_ = s.CloseWithError(masky.DefaultApplicationErrorCode, err.Error())
		}
		log.Info(fmt.Sprintf("remote client %v connect to server successfully", s.RemoteAddr()))
		go handleSession(s)
	}
}

func auth(s quic.Session) error {
	password, err := s.ReceiveMessage()
	if err != nil {
		return err
	}
	if string(password) != config.Password {
		log.Warn(fmt.Sprintf("auth error: wrong password, want: %s, get: %s", config.Password, password))
		return err
	}
	err = s.SendMessage([]byte("ok"))
	if err != nil {
		return err
	}
	return nil
}

func handleSession(s quic.Session) {
	for {
		if stream, err := s.AcceptStream(context.Background()); err == nil {
			go func() {
				if err := handleStream(&masky.Stream{Stream: stream}, s); err != nil {
					log.Error(err)
				}
			}()
		} else {
			var timeoutError *quic.IdleTimeoutError
			if !errors.As(err, &timeoutError) {
				log.Error(err)
				_ = s.CloseWithError(masky.DefaultApplicationErrorCode, err.Error())
			}
			return
		}
	}
}

func handleStream(stream *masky.Stream, s quic.Session) error {
	defer stream.Close()
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
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), addr))
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
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), req.URL.Hostname()))
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

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--port="):
			if n, err := strconv.Atoi(arg[len("--port="):]); err == nil {
				config.Port = uint16(n)
			}

		case strings.HasPrefix(arg, "--password="):
			config.Password = arg[len("--password="):]

		case strings.HasPrefix(arg, "--log="):
			switch arg[len("--log="):] {
			case "info":
				config.LogLevel = log.InfoLevel
			case "warn":
				config.LogLevel = log.WarnLevel
			case "error":
				config.LogLevel = log.ErrorLevel
			}
		}
	}
}
