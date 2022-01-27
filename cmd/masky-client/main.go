package main

import (
	"net"
	"os"
	"strings"

	"github.com/wkj9893/masky/internal/http"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"

	"github.com/wkj9893/masky/internal/socks"
)

var config masky.Config

func init() {
	// default config
	config = masky.Config{
		Port:     "2021",
		Mode:     masky.RuleMode,
		Addr:     "127.0.0.1:2022",
		LogLevel: log.InfoLevel,
	}
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
}

func main() {
	l, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	conn := masky.New(c)
	head, err := conn.Reader().Peek(1)
	if err != nil {
		return
	}
	if head[0] == 5 {
		socks.HandleConn(conn, config)
	} else {
		http.HandleConn(conn, config)
	}
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case arg == "-h", arg == "--help":

		case strings.HasPrefix(arg, "--port="):
			config.Port = arg[len("--port="):]

		case strings.HasPrefix(arg, "--mode="):
			mode := arg[len("--mode="):]
			if mode == "direct" {
				config.Mode = masky.DirectMode
			} else if mode == "global" {
				config.Mode = masky.GlobalMode
			}

		case strings.HasPrefix(arg, "--addr="):
			config.Addr = arg[len("--addr="):]

		case strings.HasPrefix(arg, "--log="):
			level := arg[len("--log="):]
			if level == "warn" {
				config.LogLevel = log.InfoLevel
			} else if level == "error" {
				config.LogLevel = log.ErrorLevel
			}
		}
	}
}
