package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/wkj9893/masky/internal/api"
	"github.com/wkj9893/masky/internal/http"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"

	"github.com/wkj9893/masky/internal/socks"
)

var config masky.ClientConfig

func init() {
	// default config
	config = masky.ClientConfig{
		Port:     "2021",
		Mode:     masky.RuleMode,
		Addr:     "127.0.0.1:2022",
		Password: "",
		AllowLan: true,
		LogLevel: log.InfoLevel,
	}
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
}

func parseAddr() string {
	if config.AllowLan {
		return ":" + config.Port
	}
	return "127.0.0.1" + ":" + config.Port
}

func main() {
	client, err := masky.NewClient(config)
	if err != nil {
		panic(err)
	}
	log.Info("connect to remote server successfully")
	l, err := net.Listen("tcp", parseAddr())
	if err != nil {
		panic(err)
	}
	log.Info("client listen on port", config.Port)
	go api.Start(client)
	log.Info("API Server listening at: http://localhost:2022")
	for {
		if c, err := l.Accept(); err == nil {
			go handleConn(c, client)
		}
	}
}

func handleConn(c net.Conn, client *masky.Client) {
	defer c.Close()
	conn := masky.NewConn(c)
	head, err := conn.Reader().Peek(1)
	if err != nil {
		log.Error(err)
		return
	}
	if head[0] == 5 {
		socks.HandleConn(conn, client)
	} else {
		http.HandleConn(conn, client)
	}
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case arg == "-h", arg == "--help":
			fmt.Println("masky client")
		case strings.HasPrefix(arg, "--port="):
			config.Port = arg[len("--port="):]

		case strings.HasPrefix(arg, "--mode="):
			mode := arg[len("--mode="):]
			if mode == "direct" {
				config.Mode = masky.DirectMode
			} else if mode == "global" {
				config.Mode = masky.GlobalMode
			} else {
				log.Warn("unknown mode option:", mode)
			}

		case strings.HasPrefix(arg, "--addr="):
			config.Addr = arg[len("--addr="):]

		case strings.HasPrefix(arg, "--dns="):
			config.Dns = arg[len("--dns="):]

		case strings.HasPrefix(arg, "--password="):
			config.Password = arg[len("--password="):]

		case strings.HasPrefix(arg, "--allowlan="):
			allowLan := arg[len("--allowlan="):]
			if allowLan == "false" {
				config.AllowLan = false
			} else if allowLan == "true" {
				config.AllowLan = true
			} else {
				log.Warn("unknown allownlan option:", allowLan)
			}

		case strings.HasPrefix(arg, "--log="):
			level := arg[len("--log="):]
			if level == "warn" {
				config.LogLevel = log.InfoLevel
			} else if level == "error" {
				config.LogLevel = log.ErrorLevel
			} else {
				log.Warn("unknown log option:", level)
			}
		}
	}
}
