package core

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/internal/logging/central"
)

// Logger wraps log/slog.Logger with OpenTelemetry trace correlation,
// TRACE and FATAL logging methods, central worker dispatch, injectable exit functions, and closer management.
type Logger struct {
	*slog.Logger
	handler  slog.Handler
	exitFunc func(int)
	closer   io.Closer
	worker   *central.Worker
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

// Worker returns the underlying central.Worker.
func (l *Logger) Worker() *central.Worker {
	if l == nil {
		return nil
	}
	return l.worker
}

// Flush blocks until all queued log events in the central worker are processed.
func (l *Logger) Flush() {
	if l != nil && l.worker != nil {
		l.worker.Flush()
	}
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
	if l == nil {
		return
	}
	l.exitFunc = resolveExitFunc(fn)
	if l.worker != nil {
		l.worker.SetExitFunc(l.exitFunc)
	}
}

// Close closes the underlying worker and closer if applicable.
func (l *Logger) Close() error {
	if l == nil {
		return nil
	}
	if err := closeWorker(l.worker); err != nil {
		return err
	}
	return closeCloser(l.closer)
}

func closeWorker(worker *central.Worker) error {
	if worker == nil {
		return nil
	}
	return worker.Close()
}

func closeCloser(closer io.Closer) error {
	if closer == nil {
		return nil
	}
	return closer.Close()
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
