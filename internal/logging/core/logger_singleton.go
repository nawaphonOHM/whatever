package core

import (
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
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
	// False positive: closing the replaced singleton is intentional ownership cleanup; proof:
	// TestGlobalLogger_SingletonResourceProof in logger_singleton_proof_test.go.
	closePreviousLogger(prev)
	l := initDefaultLogger()
	defaultLogger.Store(l)
	slog.SetDefault(l.Slog())
}
