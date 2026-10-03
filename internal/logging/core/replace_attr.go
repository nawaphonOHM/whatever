package core

import (
	"log/slog"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
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
