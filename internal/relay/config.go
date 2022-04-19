package relay

import (
	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/log"
)

type Config struct {
	Port     int       `yaml:"port"`
	LogLevel log.Level `yaml:"log_level"`
	Proxies  []Proxy   `yaml:"proxies"`
}

type Proxy struct {
	ID     uuid.UUID `yaml:"id"`
	Server string    `yaml:"server"`
}
