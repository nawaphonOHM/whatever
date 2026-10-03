package core

import (
	"io"
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
)

// NewJSON creates a JSON-formatted Logger writing to w at specified minimum level.
func NewJSON(w io.Writer, level config.Level) *Logger {
	targetW := w
	if targetW == nil {
		targetW = os.Stdout
	}
	cfg := &config.Config{
		Format: string(config.FormatJSON),
		Level:  string(level),
	}
	h := BuildHandler(cfg, targetW, nil)
	return &Logger{
		Logger:   slog.New(h),
		handler:  h,
		exitFunc: os.Exit,
		closer:   isCloserStream(targetW),
	}
}

// NewText creates a Text-formatted Logger writing to w at specified minimum level.
func NewText(w io.Writer, level config.Level) *Logger {
	targetW := w
	if targetW == nil {
		targetW = os.Stdout
	}
	cfg := &config.Config{
		Format: string(config.FormatText),
		Level:  string(level),
	}
	h := BuildHandler(cfg, targetW, nil)
	return &Logger{
		Logger:   slog.New(h),
		handler:  h,
		exitFunc: os.Exit,
		closer:   isCloserStream(targetW),
	}
}

// NewWithHandler wraps an existing slog.Handler in a Logger.
func NewWithHandler(h slog.Handler) *Logger {
	targetH := h
	if targetH == nil {
		targetH = BuildHandler(config.DefaultConfig(), os.Stdout, nil)
	}
	return &Logger{
		Logger:   slog.New(targetH),
		handler:  targetH,
		exitFunc: os.Exit,
	}
}
