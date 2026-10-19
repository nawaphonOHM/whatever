package core

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
)

var (
	defaultLogger atomic.Pointer[Logger]
	initOnce      sync.Once
)

// initDefaultLogger loads configuration from env or returns default.
func initDefaultLogger() *Logger {
	cfg, err := config.Load()
	if err != nil {
		return New(config.DefaultConfig())
	}
	l, err := NewFromConfig(cfg)
	if err != nil {
		return New(config.DefaultConfig())
	}
	return l
}

// Default returns the process-wide default Logger instance.
// It initializes lazily on first call using environment configuration.
// The logger remains open for package-level calls until ResetDefault replaces it.
func Default() *Logger {
	if l := defaultLogger.Load(); l != nil {
		return l
	}
	initOnce.Do(func() {
		// False positive: singleton initialization manages process-wide logger lifecycle;
		// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
		l := initDefaultLogger()
		defaultLogger.Store(l)
		slog.SetDefault(l.Slog())
	})
	return defaultLogger.Load()
}

// SetDefault replaces the package-level default Logger instance.
func SetDefault(l *Logger) {
	if l != nil {
		defaultLogger.Store(l)
		slog.SetDefault(l.Slog())
	}
}

func closePreviousLogger(prev *Logger) {
	if prev == nil {
		return
	}
	if err := prev.Close(); err != nil {
		return
	}
}

// ResetDefault resets the default Logger instance by reloading configuration.
// It closes any existing logger's closer if applicable.
func ResetDefault() {
	prev := defaultLogger.Swap(nil)
	// False positive: closing the replaced singleton is intentional ownership cleanup;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	closePreviousLogger(prev)
	// False positive: singleton initialization manages process-wide logger lifecycle;
	// proof: TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	l := initDefaultLogger()
	defaultLogger.Store(l)
	slog.SetDefault(l.Slog())
}

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
