package logging

import (
	"log/slog"
)

// formatLevelAttr converts slog.Level value to its canonical Level string representation.
func formatLevelAttr(a slog.Attr) slog.Attr {
	l, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}
	a.Value = slog.StringValue(string(LevelFromSlog(l)))
	return a
}

// buildReplaceAttr creates a ReplaceAttr function that standardizes level names
// and executes any custom user ReplaceAttr chain.
func buildReplaceAttr(custom ReplaceAttrFunc) ReplaceAttrFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			a = formatLevelAttr(a)
		}
		if custom != nil {
			return custom(groups, a)
		}
		return a
	}
}
