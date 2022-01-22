package masky

import (
	"github.com/wkj9893/masky/internal/log"
)

type Config struct {
	Port     string
	Mode     string // Global Rule Direct
	Addr     string
	LogLevel log.Level
}
