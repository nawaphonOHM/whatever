package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
)

// With returns a new Logger that includes the given attributes.
func (l *Logger) With(args ...any) *Logger {
	if l.isNil() {
		return nil
	}
	subSlog := l.Logger.With(args...)
	return &Logger{
		Logger:   subSlog,
		handler:  l.handler,
		exitFunc: l.exitFunc,
		closer:   l.closer,
		worker:   l.worker,
	}
}

// WithGroup returns a new Logger that starts a group.
func (l *Logger) WithGroup(name string) *Logger {
	if l.isNil() {
		return nil
	}
	subSlog := l.Logger.WithGroup(name)
	return &Logger{
		Logger:   subSlog,
		handler:  l.handler,
		exitFunc: l.exitFunc,
		closer:   l.closer,
		worker:   l.worker,
	}
}

func (l *Logger) executeExit() {
	if l.exitFunc != nil {
		l.exitFunc(1)
		return
	}
	osExit := osExitFallback()
	osExit(1)
}

func osExitFallback() func(int) {
	return resolveExitFunc(nil)
}

// Log logs at the given slog.Level through the central worker.
func (l *Logger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if l.isNil() {
		return
	}
	normCtx := normalizeContext(ctx)
	attrs := argsToAttrs(args)
	if l.worker != nil {
		l.worker.Enqueue(&central.LogEvent{
			Ctx:       normCtx,
			Logger:    l.Logger,
			Message:   msg,
			Level:     level,
			Attrs:     attrs,
			Timestamp: time.Now(),
		})
		return
	}
	l.Logger.Log(normCtx, level, msg, args...)
}

// LogAttrs logs at the given slog.Level with slog.Attr through the central worker.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if l.isNil() {
		return
	}
	normCtx := normalizeContext(ctx)
	if l.worker != nil {
		l.worker.Enqueue(&central.LogEvent{
			Ctx:       normCtx,
			Logger:    l.Logger,
			Message:   msg,
			Level:     level,
			Attrs:     attrs,
			Timestamp: time.Now(),
		})
		return
	}
	l.Logger.LogAttrs(normCtx, level, msg, attrs...)
}
