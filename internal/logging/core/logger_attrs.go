package core

import (
	"context"
	"log/slog"
)

// With returns a new Logger that includes the given attributes.
func (l *Logger) With(args ...any) *Logger {
	if l.isNil() {
		return nil
	}
	return &Logger{
		Logger:   l.Logger.With(args...),
		handler:  l.handler,
		exitFunc: l.exitFunc,
		closer:   l.closer,
	}
}

// WithGroup returns a new Logger that starts a group.
func (l *Logger) WithGroup(name string) *Logger {
	if l.isNil() {
		return nil
	}
	return &Logger{
		Logger:   l.Logger.WithGroup(name),
		handler:  l.handler,
		exitFunc: l.exitFunc,
		closer:   l.closer,
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

// Log logs at the given slog.Level.
func (l *Logger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Logger.Log(normalizeContext(ctx), level, msg, args...)
}

// LogAttrs logs at the given slog.Level with slog.Attr.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if l.isNil() {
		return
	}
	l.Logger.LogAttrs(normalizeContext(ctx), level, msg, attrs...)
}
