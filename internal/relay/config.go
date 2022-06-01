package relay

import (
	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/log"
)

type Config struct {
	Port     int       `yaml:"port"`
	Type     string    `yaml:"type"`
	Cert     string    `yaml:"cert"`
	Key      string    `yaml:"key"`
	LogLevel log.Level `yaml:"logLevel"`
	Proxies  []Proxy   `yaml:"proxies"`
}

type Proxy struct {
	ID     uuid.UUID `yaml:"id"`
	Type   string    `yaml:"type"`
	Server string    `yaml:"server"`
}
