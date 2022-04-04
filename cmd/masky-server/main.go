package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/server"
	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("c", "config.yaml", "relay config")
	flag.Parse()
	config := parseConfig(*name)
	log.SetLogLevel(config.LogLevel)
	server.Run(config)
}

func parseConfig(name string) *server.Config {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil
	}
	config := &server.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Println(err)
		return nil
	}
	return config
}
