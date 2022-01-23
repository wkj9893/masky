package log

import (
	"log"
	"os"
	"runtime/debug"
)

type Level int

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
