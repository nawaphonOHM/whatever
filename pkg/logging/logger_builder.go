package logging

import (
	"io"
	"log/slog"
	"os"
)

// resolveOutput selects non-nil writer or defaults to os.Stdout.
func resolveOutput(w io.Writer) io.Writer {
	if w != nil {
		return w
	}
	return os.Stdout
}

// resolveExitFunc selects non-nil exit function or defaults to os.Exit.
func resolveExitFunc(fn func(int)) func(int) {
	if fn != nil {
		return fn
	}
	return os.Exit
}

// buildHandlerOptions prepares slog.HandlerOptions from Config.
func buildHandlerOptions(cfg Config) *slog.HandlerOptions {
	slogLevel, err := cfg.Level.SlogLevel()
	if err != nil {
		slogLevel = SlogLevelInfo
	}
	return &slog.HandlerOptions{
		Level:       slogLevel,
		AddSource:   cfg.AddSource,
		ReplaceAttr: buildReplaceAttr(cfg.ReplaceAttr),
	}
}

// createBaseHandler constructs JSON or Text slog.Handler based on Config.
func createBaseHandler(cfg Config) slog.Handler {
	if cfg.Handler != nil {
		return cfg.Handler
	}
	out := resolveOutput(cfg.Output)
	opts := buildHandlerOptions(cfg)
	if cfg.Format == FormatText {
		return slog.NewTextHandler(out, opts)
	}
	return slog.NewJSONHandler(out, opts)
}

// buildPipelineHandler creates the full slog.Handler pipeline including trace handler.
func buildPipelineHandler(cfg Config) slog.Handler {
	base := createBaseHandler(cfg)
	if cfg.DisableTraceCorrelation {
		return base
	}
	return NewTraceHandler(base)
}

// applyFormatAndLevelDefaults fills in default level and format if blank.
func applyFormatAndLevelDefaults(cfg *Config) {
	if cfg.Level == "" {
		cfg.Level = LevelInfo
	}
	if cfg.Format == "" {
		cfg.Format = FormatJSON
	}
}

// applyOutputAndExitDefaults fills in default output and exit function if blank.
func applyOutputAndExitDefaults(cfg *Config) {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	if cfg.ExitFunc == nil {
		cfg.ExitFunc = os.Exit
	}
}

// applyConfigDefaults fills in zero-value defaults for Config.
func applyConfigDefaults(cfg *Config) {
	applyFormatAndLevelDefaults(cfg)
	applyOutputAndExitDefaults(cfg)
}

// mergeConfig applies defaults to unspecified fields in provided Config.
func mergeConfig(cfgs []Config) Config {
	if len(cfgs) == 0 {
		return DefaultConfig()
	}
	cfg := cfgs[0]
	applyConfigDefaults(&cfg)
	return cfg
}
