package client

import "github.com/wkj9893/masky/internal/log"

type Config struct {
	Port     int       `json:"port"`
	Mode     Mode      `json:"mode"`
	Addr     string    `json:"addr"`
	AllowLan bool      `json:"allowLan"`
	LogLevel log.Level `json:"logLevel"`
}
