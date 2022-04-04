package relay

import "github.com/google/uuid"

type Config struct {
	Port    uint16  `yaml:"port"`
	Proxies []Proxy `yaml:"proxies"`
}

type Proxy struct {
	ID     uuid.UUID `yaml:"id"`
	Server string    `yaml:"server"`
}
