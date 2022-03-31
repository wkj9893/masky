package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/quic"
)

// default config
var config = client.Config{
	Port:     1080,
	Mode:     client.RuleMode,
	Addr:     "127.0.0.1:3000",
	AllowLan: true,
	LogLevel: log.InfoLevel,
}

func main() {
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
	quic.SetAddr(config.Addr)
	client.Run(&config)
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
				config.Mode = client.DirectMode
			case "rule":
				config.Mode = client.RuleMode
			case "global":
				config.Mode = client.GlobalMode
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
}
