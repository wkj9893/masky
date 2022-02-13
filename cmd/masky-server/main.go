package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
		Port:     "2022",
		Password: "",
		LogLevel: log.InfoLevel,
	}
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
}

func check(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func main() {
	tlsConf, err := tls.GenerateTLSConfig()
	check(err)
	l, err := quic.ListenAddr(":"+config.Port, tlsConf, masky.DefaultQuicConfig)
	check(err)
	for {
		s, err := l.Accept(context.Background())
		if err != nil {
			continue
		}
		// auth
		password, err := s.ReceiveMessage()
		if err != nil {
			continue
		}
		if string(password) != config.Password {
			log.Warn(fmt.Sprintf("auth error: wrong password, want: %s, get: %s", config.Password, password))
			continue
		}
		err = s.SendMessage([]byte("ok"))
		if err != nil {
			continue
		}
		log.Info(fmt.Sprintf("remote client %v connect to server successfully", s.RemoteAddr()))
		go handleSession(s)
	}
}

func handleSession(s quic.Session) {
	for {
		if stream, err := s.AcceptStream(context.Background()); err == nil {
			go handleStream(&masky.Stream{Stream: stream}, s)
		}
	}
}

func handleStream(stream *masky.Stream, s quic.Session) {
	defer stream.Close()
	c := masky.NewConn(stream)

	head, err := c.Reader().Peek(1)
	if err != nil {
		log.Error(err)
		return
	}
	if head[0] == 5 { // socks
		if _, err = c.Reader().ReadByte(); err != nil {
			log.Error(err)
			return
		}
		addr, err := socks.ReadAddr(c, make([]byte, 256))
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), addr))
		dst, err := masky.Dial(addr.String())
		if err != nil {
			log.Warn(err)
			return
		}
		masky.Relay(c, dst)
	} else { // http
		req, err := http.ReadRequest(c.Reader())
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%v -> %v -> %v", s.RemoteAddr(), s.LocalAddr(), req.URL.Hostname()))
		if req.Method == http.MethodConnect {
			dst, err := masky.Dial(req.Host)
			if err != nil {
				log.Warn(err)
				return
			}
			masky.Relay(c, dst)
		} else {
			client := http.Client{Timeout: 5 * time.Second}
			req.RequestURI = ""
			resp, err := client.Do(req)
			if err != nil {
				log.Error(err)
				return
			}
			defer resp.Body.Close()
			if err = resp.Write(c); err != nil {
				log.Error(err)
			}
		}
	}
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case arg == "-h", arg == "--help":
			fmt.Println("masky server")

		case strings.HasPrefix(arg, "--password="):
			config.Password = arg[len("--password="):]

		case strings.HasPrefix(arg, "--log="):
			level := arg[len("--log="):]
			if level == "warn" {
				config.LogLevel = log.InfoLevel
			} else if level == "error" {
				config.LogLevel = log.ErrorLevel
			}

		case strings.HasPrefix(arg, "--port="):
			config.Port = arg[len("--port="):]
		}
	}
}
