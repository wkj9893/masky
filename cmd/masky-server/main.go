package main

import (
	"flag"
	"os"

	"github.com/wkj9893/masky/internal/server"
	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("c", "config.yaml", "relay config")
	flag.Parse()
	if config, err := parseConfig(*name); err != nil {
		panic(err)
	} else {
		server.Run(config)
	}
}

func parseConfig(name string) (*server.Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c := &server.Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
