package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

func Run(config *Config) {
	log.SetLogLevel(config.LogLevel)
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
	if head[0] == 5 {
		if err := handleSocks(conn, config); err != nil {
			log.Warn(err)
		}
	} else {
		if err := HandleHttp(conn, config); err != nil {
			log.Warn(err)
		}
	}
}

func ParseArgs(args []string) *Config {
	// default config
	config := &Config{
		Port:     1080,
		Mode:     RuleMode,
		Addr:     "127.0.0.1:3000",
		AllowLan: true,
		LogLevel: log.InfoLevel,
	}
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--port="):
			if n, err := strconv.Atoi(arg[len("--port="):]); err == nil {
				config.Port = uint16(n)
			}

		case strings.HasPrefix(arg, "--mode="):
			switch arg[len("--mode="):] {
			case "direct":
				config.Mode = DirectMode
			case "rule":
				config.Mode = RuleMode
			case "global":
				config.Mode = GlobalMode
			}

		case strings.HasPrefix(arg, "--addr="):
			config.Addr = arg[len("--addr="):]

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
	return config
}
