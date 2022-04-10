package main

import (
	"flag"

	"github.com/wkj9893/masky/internal/server"
)

func main() {
	name := flag.String("c", "config.yaml", "relay config")
	flag.Parse()
	if config, err := server.ParseConfig(*name); err != nil {
		panic(err)
	} else {
		server.Run(config)
	}
}
