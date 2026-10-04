package core

import (
	"io"
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
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
	return newConfiguredLogger(targetCfg, out), nil
}

func newConfiguredLogger(cfg *config.Config, out io.Writer) *Logger {
	h := BuildHandler(cfg, out, nil)
	slogL := slog.New(h)
	w := central.New(slogL, central.DefaultBufferSize, os.Exit)
	w.Start()
	return &Logger{
		Logger:   slogL,
		handler:  h,
		exitFunc: os.Exit,
		closer:   isCloserStream(out),
		worker:   w,
	}
}

func newFallbackLogger() *Logger {
	h := BuildHandler(config.DefaultConfig(), os.Stdout, nil)
	slogL := slog.New(h)
	w := central.New(slogL, central.DefaultBufferSize, os.Exit)
	w.Start()
	return &Logger{
		Logger:   slogL,
		handler:  h,
		exitFunc: os.Exit,
		worker:   w,
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
