package log

import "github.com/DilemaFixer/gog/internal/api"

type NoopLogger struct{}

func NewNoopLogger() api.Logger {
	return &NoopLogger{}
}

func (lc *NoopLogger) Log(level api.LogLevel, format string, args ...interface{}) {}
