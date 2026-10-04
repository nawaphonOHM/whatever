package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/nawaphonOHM/whatever/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/internal/logging/central"
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
	l.Log(normalizeContext(ctx), slog.LevelDebug, msg, args...)
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
	l.Log(normalizeContext(ctx), slog.LevelInfo, msg, args...)
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
	l.Log(normalizeContext(ctx), slog.LevelWarn, msg, args...)
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
	l.Log(normalizeContext(ctx), slog.LevelError, msg, args...)
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
	normCtx := normalizeContext(ctx)
	attrs := argsToAttrs(args)
	if l.worker != nil {
		l.dispatchFatal(normCtx, msg, attrs)
		return
	}
	l.Log(normCtx, config.SlogLevelFatal, msg, args...)
	l.executeExit()
}

func (l *Logger) dispatchFatal(ctx context.Context, msg string, attrs []slog.Attr) {
	l.worker.Enqueue(&central.LogEvent{
		Ctx:       ctx,
		Logger:    l.Logger,
		Message:   msg,
		Level:     config.SlogLevelFatal,
		ExitFlag:  central.ExitAbnormal,
		Stack:     callstack.FromContext(ctx),
		Attrs:     attrs,
		Timestamp: time.Now(),
	})
	l.worker.Flush()
}
