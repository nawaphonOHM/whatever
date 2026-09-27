package logging

import (
	"io"
	"log/slog"
	"os"
)

// Logger wraps log/slog.Logger with OpenTelemetry trace correlation,
// TRACE and FATAL logging methods, and injectable exit functions.
type Logger struct {
	exitFunc func(int)
	handler  slog.Handler
	*slog.Logger
}

// New creates and configures a new Logger.
func New(cfgs ...Config) *Logger {
	cfg := mergeConfig(cfgs)
	h := buildPipelineHandler(cfg)
	return &Logger{
		Logger:   slog.New(h),
		exitFunc: resolveExitFunc(cfg.ExitFunc),
		handler:  h,
	}
}

// NewJSON creates a JSON-formatted Logger writing to w at specified minimum level.
func NewJSON(w io.Writer, level Level) *Logger {
	return New(Config{
		Output: w,
		Level:  level,
		Format: FormatJSON,
	})
}

// NewText creates a Text-formatted Logger writing to w at specified minimum level.
func NewText(w io.Writer, level Level) *Logger {
	return New(Config{
		Output: w,
		Level:  level,
		Format: FormatText,
	})
}

// NewWithHandler wraps an existing slog.Handler in a Logger.
func NewWithHandler(h slog.Handler) *Logger {
	return New(Config{
		Handler: h,
	})
}

// Slog returns the underlying standard *slog.Logger.
func (l *Logger) Slog() *slog.Logger {
	if l == nil {
		return slog.Default()
	}
	return l.Logger
}

// ExitFunc returns the currently configured exit function.
func (l *Logger) ExitFunc() func(int) {
	if l == nil {
		return os.Exit
	}
	return l.exitFunc
}

// SetExitFunc updates the function invoked upon Fatal/FatalContext.
func (l *Logger) SetExitFunc(fn func(int)) {
	if l != nil {
		l.exitFunc = resolveExitFunc(fn)
	}
}
