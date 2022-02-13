package masky

import "github.com/wkj9893/masky/internal/log"

type ClientConfig struct {
	Port     string
	Mode     Mode
	Addr     string
	Dns      string
	Password string
	LogLevel log.Level
}

type ServerConfig struct {
	Port     string
	Password string
	LogLevel log.Level
}
