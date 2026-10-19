package core

import (
	"io"
	"log/slog"
	"os"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
)

// ReplaceAttrFunc specifies a function to modify attributes before logging.
type ReplaceAttrFunc = func([]string, slog.Attr) slog.Attr

// FormatLevelAttr converts slog.Level value to its canonical Level string representation.
func FormatLevelAttr(a slog.Attr) slog.Attr {
	l, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}
	a.Value = slog.StringValue(string(config.LevelFromSlog(l)))
	return a
}

// BuildReplaceAttr creates a ReplaceAttr function that standardizes level names
// and executes any custom user ReplaceAttr chain.
func BuildReplaceAttr(custom ReplaceAttrFunc) ReplaceAttrFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			a = FormatLevelAttr(a)
		}
		if custom != nil {
			return custom(groups, a)
		}
		return a
	}
}

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
