package log

import (
	"encoding/json"
	"errors"
	"log"
	"os"
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
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ltime)
	warnLogger = log.New(os.Stdout, "WARNING: ", log.Ltime)
	errorLogger = log.New(os.Stdout, "ERROR: ", log.Ltime)
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
	}
}

func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l *Level) UnmarshalJSON(data []byte) error {
	var level string
	err := json.Unmarshal(data, &level)
	if err != nil {
		return err
	}
	switch level {
	case "info":
		*l = InfoLevel
	case "warn":
		*l = WarnLevel
	case "error":
		*l = ErrorLevel
	default:
		return errors.New("unknown logLevel")
	}
	return nil
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
