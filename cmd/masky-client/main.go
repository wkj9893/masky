package main

import (
	"flag"

	"github.com/wkj9893/masky/internal/client"
)

func main() {
	name := flag.String("c", "config.yaml", "client config")
	flag.Parse()
	if config, err := client.ParseConfig(*name); err != nil {
		panic(err)
	} else {
		client.Run(config)
	}
}
