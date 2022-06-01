package main

import (
	"flag"
	"os"

	"github.com/wkj9893/masky/internal/exit"
	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("c", "config.yaml", "relay config")
	flag.Parse()
	if config, err := parseConfig(*name); err != nil {
		panic(err)
	} else {
		exit.Run(config)
	}
}

func parseConfig(name string) (*exit.Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c := &exit.Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
