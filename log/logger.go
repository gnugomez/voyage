package log

import "github.com/charmbracelet/log"

type Logger interface {
	Info(message string, vals ...any)
	Error(message string, vals ...any)
	Fatal(message string, vals ...any)
	Debug(message string, vals ...any)
}

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

var GlobalLogger Logger = CreateDefaultLogger(InfoLevel)

func SetLogger(logger Logger) {
	GlobalLogger = logger
}

func Info(message string, vals ...any) {
	GlobalLogger.Info(message, vals...)
}

func Error(message string, vals ...any) {
	GlobalLogger.Error(message, vals...)
}

func Fatal(message string, vals ...any) {
	GlobalLogger.Fatal(message, vals...)
}

func Debug(message string, vals ...any) {
	GlobalLogger.Debug(message, vals...)
}

// Default logger implementation

var DefaultLogLevel = map[LogLevel]log.Level{
	DebugLevel: log.DebugLevel,
	InfoLevel:  log.InfoLevel,
	ErrorLevel: log.ErrorLevel,
	FatalLevel: log.FatalLevel,
}

type defaultLogger struct {
}

func CreateDefaultLogger(level LogLevel) *defaultLogger {
	log.SetLevel(DefaultLogLevel[level])
	return &defaultLogger{}
}

func (*defaultLogger) Info(message string, vals ...any) {
	log.Info(message, vals...)
}

func (*defaultLogger) Debug(message string, vals ...any) {
	log.Debug(message, vals...)
}

func (*defaultLogger) Error(message string, vals ...any) {
	log.Error(message, vals...)
}

func (*defaultLogger) Fatal(message string, vals ...any) {
	log.Fatal(message, vals...)
}
