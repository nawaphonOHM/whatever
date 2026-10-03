package core

import (
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
)

// resolveExitFunc selects non-nil exit function or defaults to os.Exit.
func resolveExitFunc(fn func(int)) func(int) {
	if fn != nil {
		return fn
	}
	return os.Exit
}

// NewFromConfig creates a new Logger configured from the provided Config.
// If cfg is nil, default configuration is used.
func NewFromConfig(cfg *config.Config) (*Logger, error) {
	targetCfg := cfg
	if targetCfg == nil {
		targetCfg = config.DefaultConfig()
	}
	out, err := targetCfg.ResolveOutput()
	if err != nil {
		return nil, err
	}
	h := BuildHandler(targetCfg, out, nil)
	return &Logger{
		Logger:   slog.New(h),
		handler:  h,
		exitFunc: os.Exit,
		closer:   isCloserStream(out),
	}, nil
}

func newFallbackLogger() *Logger {
	h := BuildHandler(config.DefaultConfig(), os.Stdout, nil)
	return &Logger{
		Logger:   slog.New(h),
		handler:  h,
		exitFunc: os.Exit,
	}
}

func resolveConfig(cfgs []*config.Config) *config.Config {
	if len(cfgs) > 0 && cfgs[0] != nil {
		return cfgs[0]
	}
	return config.DefaultConfig()
}

// New creates a new Logger from optional config.Config pointer.
// If no config is provided or error occurs, it falls back to defaults.
func New(cfgs ...*config.Config) *Logger {
	l, err := NewFromConfig(resolveConfig(cfgs))
	if err != nil {
		return newFallbackLogger()
	}
	return l
}
