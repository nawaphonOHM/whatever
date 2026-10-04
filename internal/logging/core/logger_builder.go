package core

import (
	"io"
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
)

func isStandardIO(w io.Writer) bool {
	return w == nil || w == os.Stdout || w == os.Stderr
}

func isDiscardOrStdin(w io.Writer) bool {
	return w == os.Stdin || w == io.Discard
}

func isSpecialStream(w io.Writer) bool {
	if isStandardIO(w) {
		return true
	}
	return isDiscardOrStdin(w)
}

func isCloserStream(w io.Writer) io.Closer {
	if isSpecialStream(w) {
		return nil
	}
	if c, ok := w.(io.Closer); ok {
		return c
	}
	return nil
}

// BuildHandlerOptions prepares slog.HandlerOptions from config.Config.
func BuildHandlerOptions(cfg *config.Config, customReplaceAttr ReplaceAttrFunc) *slog.HandlerOptions {
	slogLevel := config.SlogLevelInfo
	var addSource bool
	if cfg != nil {
		lvl, err := cfg.ParsedSlogLevel()
		if err == nil {
			slogLevel = lvl
		}
		addSource = cfg.AddSource
	}
	return &slog.HandlerOptions{
		Level:       slogLevel,
		AddSource:   addSource,
		ReplaceAttr: BuildReplaceAttr(customReplaceAttr),
	}
}

func resolveBaseFormat(cfg *config.Config) config.Format {
	if cfg == nil {
		return config.FormatText
	}
	f, err := cfg.ParsedFormat()
	if err != nil {
		return config.FormatText
	}
	return f
}

// BuildBaseHandler constructs JSON or Text slog.Handler based on Config and output writer.
func BuildBaseHandler(cfg *config.Config, out io.Writer, customReplaceAttr ReplaceAttrFunc) slog.Handler {
	targetOut := out
	if targetOut == nil {
		targetOut = os.Stdout
	}
	opts := BuildHandlerOptions(cfg, customReplaceAttr)
	if resolveBaseFormat(cfg) == config.FormatJSON {
		return slog.NewJSONHandler(targetOut, opts)
	}
	return slog.NewTextHandler(targetOut, opts)
}

// BuildHandler creates the full slog.Handler pipeline including trace handler according to Config.
func BuildHandler(cfg *config.Config, out io.Writer, customReplaceAttr ReplaceAttrFunc) slog.Handler {
	base := BuildBaseHandler(cfg, out, customReplaceAttr)
	if cfg != nil && cfg.DisableTraceCorrelation {
		return base
	}
	return NewTraceHandler(base)
}
