package core

import (
	"context"
	"log/slog"
)

// False positive: package-level logging intentionally borrows the process-wide logger instead of closing it per call;
// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.

// Trace logs at TRACE level using the default logger.
func Trace(msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Trace(msg, args...)
}

// TraceContext logs at TRACE level with context using the default logger.
func TraceContext(ctx context.Context, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().TraceContext(ctx, msg, args...)
}

// Debug logs at DEBUG level using the default logger.
func Debug(msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Debug(msg, args...)
}

// DebugContext logs at DEBUG level with context using the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().DebugContext(ctx, msg, args...)
}

// Info logs at INFO level using the default logger.
func Info(msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Info(msg, args...)
}

// InfoContext logs at INFO level with context using the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().InfoContext(ctx, msg, args...)
}

// Warn logs at WARN level using the default logger.
func Warn(msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Warn(msg, args...)
}

// WarnContext logs at WARN level with context using the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().WarnContext(ctx, msg, args...)
}

// Error logs at ERROR level using the default logger.
func Error(msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Error(msg, args...)
}

// ErrorContext logs at ERROR level with context using the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().ErrorContext(ctx, msg, args...)
}

// Fatal logs at FATAL level and exits using the default logger.
func Fatal(msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Fatal(msg, args...)
}

// FatalContext logs at FATAL level with context and exits using the default logger.
func FatalContext(ctx context.Context, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().FatalContext(ctx, msg, args...)
}

// Log logs at the given level using the default logger.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Log(ctx, level, msg, args...)
}

// LogAttrs logs at the given level with attributes using the default logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().LogAttrs(ctx, level, msg, attrs...)
}

// Flush flushes all buffered events in the default logger's central worker.
func Flush() {
	// False positive: this wrapper borrows the singleton;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	Default().Flush()
}
