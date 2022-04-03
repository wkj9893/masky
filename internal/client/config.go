package client

import (
	"github.com/wkj9893/masky/internal/log"
)

type Config struct {
	Port     uint16    `yaml:"port"`
	Mode     Mode      `yaml:"mode"`
	AllowLan bool      `yaml:"allow_lan"`
	LogLevel log.Level `yaml:"log_level"`
	Proxies  []Proxy   `yaml:"proxies"`
}

type Proxy struct {
	Name   string   `yaml:"name"`
	Server []string `yaml:"server"`
}
