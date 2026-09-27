package logging

import (
	"context"
	"log/slog"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(New())
}

// Default returns the package-level default Logger instance.
func Default() *Logger {
	return defaultLogger.Load()
}

// SetDefault replaces the package-level default Logger instance.
func SetDefault(l *Logger) {
	if l != nil {
		defaultLogger.Store(l)
		slog.SetDefault(l.Logger)
	}
}

// Trace logs a message at TRACE level using the default logger.
func Trace(msg string, args ...any) {
	Default().Trace(msg, args...)
}

// TraceContext logs a message at TRACE level with context using the default logger.
func TraceContext(ctx context.Context, msg string, args ...any) {
	Default().TraceContext(ctx, msg, args...)
}

// Debug logs a message at DEBUG level using the default logger.
func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

// DebugContext logs a message at DEBUG level with context using the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Default().DebugContext(ctx, msg, args...)
}

// Info logs a message at INFO level using the default logger.
func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

// InfoContext logs a message at INFO level with context using the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Default().InfoContext(ctx, msg, args...)
}

// Warn logs a message at WARN level using the default logger.
func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

// WarnContext logs a message at WARN level with context using the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	Default().WarnContext(ctx, msg, args...)
}

// Error logs a message at ERROR level using the default logger.
func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

// ErrorContext logs a message at ERROR level with context using the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Default().ErrorContext(ctx, msg, args...)
}

// Fatal logs a message at FATAL level and exits using the default logger.
func Fatal(msg string, args ...any) {
	Default().Fatal(msg, args...)
}

// FatalContext logs a message at FATAL level with context and exits using the default logger.
func FatalContext(ctx context.Context, msg string, args ...any) {
	Default().FatalContext(ctx, msg, args...)
}
