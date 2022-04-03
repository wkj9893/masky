package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wkj9893/masky/internal/client"
	"github.com/wkj9893/masky/internal/log"
	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("c", "config.yaml", "client config")
	flag.Parse()
	config := parseConfig(*name)
	log.SetLogLevel(config.LogLevel)
	client.Run(config)
}

func parseConfig(name string) *client.Config {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil
	}
	config := &client.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Println(err)
		return nil
	}
	return config
}
