package client

import (
	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/log"
)

type Config struct {
	Port     int       `yaml:"port" json:"port"`
	Mode     Mode      `yaml:"mode" json:"mode"`
	AllowLan bool      `yaml:"allowLan" json:"allowLan"`
	LogLevel log.Level `yaml:"logLevel" json:"logLevel"`
	Proxies  []Proxy   `yaml:"proxies" json:"proxies"`
}

type Proxy struct {
	ID     uuid.UUID `yaml:"id" json:"id"`
	Name   string    `yaml:"name" json:"name"`
	Type   string    `yaml:"type" json:"type"`
	Server []string  `yaml:"server" json:"server"`
}
