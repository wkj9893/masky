package main

import (
	"flag"

	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/log"
)

func main() {
	client.Run(parseArgs())
}

func parseArgs() *client.Config {
	port := flag.Int("port", 1080, "local listen port")
	mode := flag.String("mode", "rule", "client mode(direct|rule|global)")
	addr := flag.String("addr", "127.0.0.1:1081", "remote server addr")
	allowLan := flag.Bool("lan", false, "allow lan")
	logLevel := flag.String("log", "info", "logLevel(info|warn|error)")
	flag.Parse()
	return &client.Config{
		Port:     *port,
		Mode:     client.NewMode(*mode),
		Addr:     *addr,
		AllowLan: *allowLan,
		LogLevel: log.NewLevel(*logLevel),
	}
}
