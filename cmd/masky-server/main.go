package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/server"
)

// default config
var config = server.Config{
	Port:     3000,
	LogLevel: log.InfoLevel,
}

func main() {
	parseArgs(os.Args[1:])
	log.SetLogLevel(config.LogLevel)
	server.Run(&config)
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--port="):
			if n, err := strconv.Atoi(arg[len("--port="):]); err == nil {
				config.Port = uint16(n)
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
