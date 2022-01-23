package masky

import (
	"github.com/wkj9893/masky/internal/log"
)

type Config struct {
	Port     string
	Mode     Mode
	Addr     string
	LogLevel log.Level
}
