package main

import (
	"flag"
	"os"

	"github.com/wkj9893/masky/internal/client"
	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("c", "config.yaml", "client config")
	flag.Parse()
	if config, err := parseConfig(*name); err != nil {
		panic(err)
	} else {
		client.Run(config)
	}
}

func parseConfig(name string) (*client.Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c := &client.Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
