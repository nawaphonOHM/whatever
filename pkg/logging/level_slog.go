package logging

import (
	"log/slog"
)

// mapLowerSlogLevel maps negative slog levels to TRACE or DEBUG.
func mapLowerSlogLevel(l slog.Level) (Level, bool) {
	if l <= SlogLevelTrace {
		return LevelTrace, true
	}
	if l < SlogLevelInfo {
		return LevelDebug, true
	}
	return "", false
}

// mapUpperSlogLevel maps positive slog levels to WARN, ERROR, or FATAL.
func mapUpperSlogLevel(l slog.Level) Level {
	if l < SlogLevelError {
		return LevelWarn
	}
	if l < SlogLevelFatal {
		return LevelError
	}
	return LevelFatal
}

// LevelFromSlog maps any slog.Level to its canonical Level representation.
func LevelFromSlog(l slog.Level) Level {
	if lvl, ok := mapLowerSlogLevel(l); ok {
		return lvl
	}
	if l < SlogLevelWarn {
		return LevelInfo
	}
	return mapUpperSlogLevel(l)
}
