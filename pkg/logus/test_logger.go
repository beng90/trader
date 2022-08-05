package logus

type TestLogger struct {
}

func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (l *TestLogger) Debug(v ...any) {
}

func (l *TestLogger) Debugf(format string, v ...any) {
}

func (l *TestLogger) Error(v ...any) {
}

func (l *TestLogger) Fatal(v ...any) {
}
