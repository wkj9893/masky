package log

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
)

type Level uint8

const (
	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger

	logLevel Level
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", 0)
	warnLogger = log.New(os.Stdout, "WARNING: ", 0)
	errorLogger = log.New(os.Stdout, "ERROR: ", 0)
}

func SetLogLevel(l Level) {
	logLevel = l
}

func Info(v ...interface{}) {
	if logLevel <= InfoLevel {
		infoLogger.Println(v...)
	}
}

func Warn(v ...interface{}) {
	if logLevel <= WarnLevel {
		warnLogger.Println(v...)
	}
}

func Error(v ...interface{}) {
	if logLevel <= ErrorLevel {
		errorLogger.Println(v...)
		debug.PrintStack()
	}
}

func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l Level) String() string {
	switch l {
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "unknown"
	}
}
