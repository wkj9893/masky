package log

import (
	"log"
	"os"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	SilentLevel
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

func Info(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		infoLogger.Printf(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= WarnLevel {
		warnLogger.Printf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ErrorLevel {
		errorLogger.Printf(format, v...)
	}
}
