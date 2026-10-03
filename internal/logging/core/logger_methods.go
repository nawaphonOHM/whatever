package core

import (
	"context"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
)

// Trace logs at TRACE level.
func (l *Logger) Trace(msg string, args ...any) {
	l.TraceContext(context.Background(), msg, args...)
}

// TraceContext logs at TRACE level with context.
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Log(normalizeContext(ctx), config.SlogLevelTrace, msg, args...)
}

// Debug logs at DEBUG level.
func (l *Logger) Debug(msg string, args ...any) {
	l.DebugContext(context.Background(), msg, args...)
}

// DebugContext logs at DEBUG level with context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Logger.DebugContext(normalizeContext(ctx), msg, args...)
}

// Info logs at INFO level.
func (l *Logger) Info(msg string, args ...any) {
	l.InfoContext(context.Background(), msg, args...)
}

// InfoContext logs at INFO level with context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Logger.InfoContext(normalizeContext(ctx), msg, args...)
}

// Warn logs at WARN level.
func (l *Logger) Warn(msg string, args ...any) {
	l.WarnContext(context.Background(), msg, args...)
}

// WarnContext logs at WARN level with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Logger.WarnContext(normalizeContext(ctx), msg, args...)
}

// Error logs at ERROR level.
func (l *Logger) Error(msg string, args ...any) {
	l.ErrorContext(context.Background(), msg, args...)
}

// ErrorContext logs at ERROR level with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Logger.ErrorContext(normalizeContext(ctx), msg, args...)
}

// Fatal logs at FATAL level and exits.
func (l *Logger) Fatal(msg string, args ...any) {
	l.FatalContext(context.Background(), msg, args...)
}

// FatalContext logs at FATAL level with context and exits.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Log(normalizeContext(ctx), config.SlogLevelFatal, msg, args...)
	l.executeExit()
}
