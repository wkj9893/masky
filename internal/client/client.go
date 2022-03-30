package client

import (
	"fmt"
	"net"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func Run(config Config) {
	addr := fmt.Sprintf(":%v", config.Port)
	if !config.AllowLan {
		addr = fmt.Sprintf("127.0.0.1:%v", config.Port)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Info("client listen on port", config.Port)
	for {
		if c, err := l.Accept(); err == nil {
			go handleConn(c, config.Mode)
		} else {
			panic(err)
		}
	}
}

func handleConn(c net.Conn, mode Mode) {
	defer c.Close()
	conn := masky.NewConn(c)
	head, err := conn.Reader().Peek(1)
	if err != nil {
		log.Error(err)
		return
	}
	if head[0] == 5 {
		if err := handleSocks(conn, mode); err != nil {
			log.Warn(err)
		}
	} else {
		if err := HandleHttp(conn, mode); err != nil {
			log.Warn(err)
		}
	}
}
