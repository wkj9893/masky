package main

import (
	"flag"
	"net"
	"strconv"

	h "github.com/wkj9893/masky/internal/http"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
	"github.com/wkj9893/masky/internal/socks"
)

var config masky.Config

func init() {
	flag.IntVar(&config.Port, "port", 2021, "Listen Port")
	flag.StringVar(&config.Mode, "mode", "Rule", "Client Mode")
	flag.BoolVar(&config.Debug, "debug", false, "Debug")
	flag.StringVar(&config.Addr, "addr", ":2022", "Server Address")
	flag.Parse()
}

func main() {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Error("%v", err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Error("%v", err)
		}
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	conn := masky.New(c)
	head, err := conn.Reader().Peek(1)
	if err != nil {
		log.Error("%v", err)
		return
	}
	if head[0] == 5 {
		socks.HandleConn(conn, config)
	} else {
		h.HandleConn(conn, config)
	}
}
