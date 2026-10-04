package logging

import (
	"context"
	"log/slog"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/nawaphonOHM/whatever/internal/logging/core"
)

// Trace logs a message at TRACE level using the default logger.
func Trace(msg string, args ...any) {
	core.Trace(msg, args...)
}

// TraceContext logs a message at TRACE level with context using the default logger.
func TraceContext(ctx context.Context, msg string, args ...any) {
	core.TraceContext(ctx, msg, args...)
}

// Debug logs a message at DEBUG level using the default logger.
func Debug(msg string, args ...any) {
	core.Debug(msg, args...)
}

// DebugContext logs a message at DEBUG level with context using the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	core.DebugContext(ctx, msg, args...)
}

// Info logs a message at INFO level using the default logger.
func Info(msg string, args ...any) {
	core.Info(msg, args...)
}

// InfoContext logs a message at INFO level with context using the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	core.InfoContext(ctx, msg, args...)
}

// Warn logs a message at WARN level using the default logger.
func Warn(msg string, args ...any) {
	core.Warn(msg, args...)
}

// WarnContext logs a message at WARN level with context using the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	core.WarnContext(ctx, msg, args...)
}

// Error logs a message at ERROR level using the default logger.
func Error(msg string, args ...any) {
	core.Error(msg, args...)
}

// ErrorContext logs a message at ERROR level with context using the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	core.ErrorContext(ctx, msg, args...)
}

// Fatal logs a message at FATAL level and exits using the default logger.
func Fatal(msg string, args ...any) {
	core.Fatal(msg, args...)
}

// FatalContext logs a message at FATAL level with context and exits using the default logger.
func FatalContext(ctx context.Context, msg string, args ...any) {
	core.FatalContext(ctx, msg, args...)
}

// Log logs a message at the given level using the default logger.
func Log(ctx context.Context, level Level, msg string, args ...any) {
	slogLevel, err := level.SlogLevel()
	if err != nil {
		slogLevel = config.SlogLevelInfo
	}
	core.Log(ctx, slogLevel, msg, args...)
}

// LogAttrs logs a message at the given level with attributes using the default logger.
func LogAttrs(ctx context.Context, level Level, msg string, attrs ...slog.Attr) {
	slogLevel, err := level.SlogLevel()
	if err != nil {
		slogLevel = config.SlogLevelInfo
	}
	core.LogAttrs(ctx, slogLevel, msg, attrs...)
}

// Flush flushes any buffered log entries in the default central logger.
func Flush() {
	core.Flush()
}
