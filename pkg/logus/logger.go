package logus

import (
	"log"

	"gorm.io/gorm/logger"
)

type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

type StdLogger struct {
	logger   *log.Logger
	logLevel logger.LogLevel
}

func NewStdLogger(logger *log.Logger) *StdLogger {
	return &StdLogger{logger: logger}
}

func (l *StdLogger) Debug(v ...any) {
	if l.logLevel != logger.Info {
		return
	}

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
