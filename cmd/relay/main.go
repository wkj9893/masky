package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wkj9893/masky/internal/relay"
	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("c", "config.yaml", "relay config")
	flag.Parse()
	config := parseConfig(*name)
	relay.Run(config)
}

func parseConfig(name string) *relay.Config {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil
	}
	config := &relay.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Println(err)
		return nil
	}
	return config
}
