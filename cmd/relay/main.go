package main

import (
	"flag"

	"github.com/wkj9893/masky/internal/relay"
)

func main() {
	name := flag.String("c", "config.yaml", "relay config")
	flag.Parse()
	if config, err := relay.ParseConfig(*name); err != nil {
		panic(err)
	} else {
		relay.Run(config)
	}
}
