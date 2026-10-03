package core

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Logger wraps log/slog.Logger with OpenTelemetry trace correlation,
// TRACE and FATAL logging methods, injectable exit functions, and closer management.
type Logger struct {
	*slog.Logger
	handler  slog.Handler
	exitFunc func(int)
	closer   io.Closer
}

// Slog returns the underlying standard *slog.Logger.
func (l *Logger) Slog() *slog.Logger {
	if l == nil || l.Logger == nil {
		return slog.Default()
	}
	return l.Logger
}

// Handler returns the underlying slog.Handler.
func (l *Logger) Handler() slog.Handler {
	if l == nil {
		return nil
	}
	return l.handler
}

// ExitFunc returns the currently configured exit function.
func (l *Logger) ExitFunc() func(int) {
	if l == nil || l.exitFunc == nil {
		return os.Exit
	}
	return l.exitFunc
}

// SetExitFunc updates the function invoked upon Fatal/FatalContext.
func (l *Logger) SetExitFunc(fn func(int)) {
	if l != nil {
		l.exitFunc = resolveExitFunc(fn)
	}
}

// Close closes the underlying resource if it implements io.Closer.
func (l *Logger) Close() error {
	if l == nil || l.closer == nil {
		return nil
	}
	return l.closer.Close()
}

func (l *Logger) isNil() bool {
	return l == nil || l.Logger == nil
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
