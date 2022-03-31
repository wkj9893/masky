package client

import "github.com/wkj9893/masky/internal/log"

type Config struct {
	Port     uint16    `json:"port"`
	Mode     Mode      `json:"mode"`
	Addrs    []string  `json:"addr"`
	AllowLan bool      `json:"allowLan"`
	LogLevel log.Level `json:"logLevel"`
}
