package server

import "github.com/wkj9893/masky/internal/log"

type Config struct {
	Port     uint16
	LogLevel log.Level
}
