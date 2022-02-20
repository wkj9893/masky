package masky

import (
	"github.com/wkj9893/masky/internal/log"
)

type ClientConfig struct {
	Port     string    `json:"port"`
	Mode     Mode      `json:"mode"`
	Addr     string    `json:"addr"`
	Dns      string    `json:"dns"`
	Password string    `json:"password"`
	AllowLan bool      `json:"allow_lan"`
	LogLevel log.Level `json:"log_level"`
}

type ServerConfig struct {
	Port     string
	Password string
	LogLevel log.Level
}
