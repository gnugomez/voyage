package log

import "github.com/charmbracelet/log"

var DefaultLogLevel = map[LogLevel]log.Level{
	DebugLevel: log.DebugLevel,
	InfoLevel:  log.InfoLevel,
	ErrorLevel: log.ErrorLevel,
	FatalLevel: log.FatalLevel,
}

type DefaultLogger struct {
}

func CreateDefaultLogger(level LogLevel) *DefaultLogger {
	log.SetLevel(DefaultLogLevel[level])
	return &DefaultLogger{}
}

func (*DefaultLogger) Info(message string, vals ...any) {
	log.Info(message, vals...)
}

func (*DefaultLogger) Debug(message string, vals ...any) {
	log.Debug(message, vals...)
}

func (*DefaultLogger) Error(message string, vals ...any) {
	log.Error(message, vals)
}

func (*DefaultLogger) Fatal(message string, vals ...any) {
	log.Fatal(message, vals...)
}
