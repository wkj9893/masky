package client

import (
	"fmt"
	"net"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func Run(config *Config) {
	log.SetLogLevel(config.LogLevel)
	addr := fmt.Sprintf("127.0.0.1:%v", config.Port)
	if !config.AllowLan {
		addr = fmt.Sprintf(":%v", config.Port)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Info("client listen on port", config.Port)
	for {
		if c, err := l.Accept(); err == nil {
			go handleConn(c, config)
		} else {
			panic(err)
		}
	}
}

func handleConn(c net.Conn, config *Config) {
	defer c.Close()
	conn := masky.NewConn(c)
	head, err := conn.Reader().Peek(1)
	if err != nil {
		log.Error(err)
		return
	}
	switch head[0] {
	case 5: // socks
		if err := handleSocks(conn, config); err != nil {
			log.Warn(err)
		}
	default: // http
		if err := handleHttp(conn, config); err != nil {
			log.Warn(err)
		}
	}
}
