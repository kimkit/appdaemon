package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	LogDebug(string, string, ...interface{})
	LogInfo(string, string, ...interface{})
	LogWarning(string, string, ...interface{})
	LogError(string, string, ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) LogDebug(prefix, format string, v ...interface{}) {
	log.Printf("DEBUG %s: %s", prefix, fmt.Sprintf(format, v...))
}

func (l *defaultLogger) LogInfo(prefix, format string, v ...interface{}) {
	log.Printf("INFO %s: %s", prefix, fmt.Sprintf(format, v...))
}

func (l *defaultLogger) LogWarning(prefix, format string, v ...interface{}) {
	log.Printf("WARNING %s: %s", prefix, fmt.Sprintf(format, v...))
}

func (l *defaultLogger) LogError(prefix, format string, v ...interface{}) {
	log.Printf("ERROR %s: %s", prefix, fmt.Sprintf(format, v...))
}

func NewLogger() Logger {
	return &defaultLogger{}
}
