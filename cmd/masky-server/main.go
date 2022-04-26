package main

import (
	"flag"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/server"
)

func main() {
	server.Run(parseArgs())
}

func parseArgs() *server.Config {
	port := flag.Int("port", 1081, "local listen port")
	logLevel := flag.String("log", "info", "logLevel(info|warn|error)")
	flag.Parse()
	return &server.Config{
		Port:     *port,
		LogLevel: log.NewLevel(*logLevel),
	}
}
