package core

import (
	"io"
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/internal/logging/central"
	"github.com/nawaphonOHM/whatever/internal/logging/config"
)

// NewJSON creates a JSON-formatted Logger writing to w at specified minimum level.
func NewJSON(w io.Writer, level config.Level) *Logger {
	return newWriterLogger(w, level, config.FormatJSON)
}

func newWriterLogger(w io.Writer, level config.Level, format config.Format) *Logger {
	targetW := resolveWriter(w)
	cfg := &config.Config{
		Format: string(format),
		Level:  string(level),
	}
	h := BuildHandler(cfg, targetW, nil)
	slogL := slog.New(h)
	return buildWriterLogger(slogL, h, targetW)
}

func resolveWriter(w io.Writer) io.Writer {
	if w == nil {
		return os.Stdout
	}
	return w
}

func buildWriterLogger(slogL *slog.Logger, h slog.Handler, targetW io.Writer) *Logger {
	bufSize := central.DefaultBufferSize
	if !isSpecialStream(targetW) {
		bufSize = 0
	}
	worker := central.New(slogL, bufSize, os.Exit)
	worker.Start()
	return &Logger{
		Logger:   slogL,
		handler:  h,
		exitFunc: os.Exit,
		closer:   isCloserStream(targetW),
		worker:   worker,
	}
}

// NewText creates a Text-formatted Logger writing to w at specified minimum level.
func NewText(w io.Writer, level config.Level) *Logger {
	return newWriterLogger(w, level, config.FormatText)
}

// NewWithHandler wraps an existing slog.Handler in a Logger.
func NewWithHandler(h slog.Handler) *Logger {
	targetH := h
	if targetH == nil {
		targetH = BuildHandler(config.DefaultConfig(), os.Stdout, nil)
	}
	slogL := slog.New(targetH)
	worker := central.New(slogL, 0, os.Exit)
	worker.Start()
	return &Logger{
		Logger:   slogL,
		handler:  targetH,
		exitFunc: os.Exit,
		worker:   worker,
	}
}
