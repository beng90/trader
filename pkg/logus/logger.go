package logus

import (
	"log"
)

type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

type StdLogger struct {
	logger *log.Logger
}

func NewStdLogger(logger *log.Logger) *StdLogger {
	return &StdLogger{logger: logger}
}

func (l *StdLogger) Debug(v ...any) {
	l.logger.Println(v)
}

func (l *StdLogger) Debugf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l *StdLogger) Error(v ...any) {
	l.logger.Println(v)
}

func (l *StdLogger) Fatal(v ...any) {
	l.logger.Fatalln(v)
}
