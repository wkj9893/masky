package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/wkj9893/masky/internal/relay"
)

var config = relay.Config{
	Port:  3001,
	Addrs: []string{},
}

func main() {
	parseArgs(os.Args[1:])
	relay.Run(&config)
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--port="):
			if n, err := strconv.Atoi(arg[len("--port="):]); err == nil {
				config.Port = uint16(n)
			}

		case strings.HasPrefix(arg, "--addrs="):
			addrs := strings.Split(arg[len("--addr="):], ",")
			for _, addr := range addrs {
				config.Addrs = append(config.Addrs, strings.TrimSpace(addr))
			}
		}
	}
}
