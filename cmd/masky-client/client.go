package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/wkj9893/masky/internal/api"
	"github.com/wkj9893/masky/internal/http"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"

	"github.com/wkj9893/masky/internal/socks"
)

var (
	config    masky.ClientConfig
	localAddr string
)

func init() {
	// default config
	config = masky.ClientConfig{
		Port:     1080,
		Mode:     masky.RuleMode,
		Addr:     "127.0.0.1:1080",
		Dns:      "",
		Password: "masky",
		AllowLan: true,
		LogLevel: log.InfoLevel,
	}
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
	if config.AllowLan {
		localAddr = fmt.Sprintf(":%v", config.Port)
	} else {
		localAddr = fmt.Sprintf("127.0.0.1:%v", config.Port)
	}
}

func main() {
	client, err := masky.NewClient(config)
	if err != nil {
		panic(err)
	}
	log.Info("connect to remote server successfully")
	l, err := net.Listen("tcp", localAddr)
	if err != nil {
		panic(err)
	}
	log.Info("client listen on port", config.Port)
	go func() {
		if err := api.Start(client); err != nil {
			log.Warn("fail to start api server", err)
		}
	}()
	log.Info("API Server listening at: http://localhost:1081")
	for {
		if c, err := l.Accept(); err == nil {
			go handleConn(c, client)
		} else {
			panic(err)
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
		if err := socks.HandleConn(conn, client); err != nil {
			log.Warn(err)
		}
	} else {
		if err := http.HandleConn(conn, client); err != nil {
			log.Warn(err)
		}
	}
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--port="):
			if n, err := strconv.Atoi(arg[len("--port="):]); err == nil {
				config.Port = uint16(n)
			}

		case strings.HasPrefix(arg, "--mode="):
			switch arg[len("--mode="):] {
			case "direct":
				config.Mode = masky.DirectMode
			case "rule":
				config.Mode = masky.RuleMode
			case "global":
				config.Mode = masky.GlobalMode
			}

		case strings.HasPrefix(arg, "--addr="):
			config.Addr = arg[len("--addr="):]

		case strings.HasPrefix(arg, "--dns="):
			config.Dns = arg[len("--dns="):]

		case strings.HasPrefix(arg, "--password="):
			config.Password = arg[len("--password="):]

		case strings.HasPrefix(arg, "--allowlan="):
			switch arg[len("--allowlan="):] {
			case "true":
				config.AllowLan = true
			case "false":
				config.AllowLan = false
			}

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
