package masky

import (
	"github.com/wkj9893/masky/internal/log"
)

type ClientConfig struct {
	Port     uint16    `json:"port"`
	Mode     Mode      `json:"mode"`
	Addr     string    `json:"addr"`
	Password string    `json:"password"`
	AllowLan bool      `json:"allowLan"`
	LogLevel log.Level `json:"logLevel"`
}

type ServerConfig struct {
	Port     uint16
	Password string
	LogLevel log.Level
}
