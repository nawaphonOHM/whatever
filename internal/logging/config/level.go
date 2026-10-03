package config

import (
	"log/slog"
)

// Level represents supported logging severity levels.
type Level string

// Supported log level string constants.
const (
	LevelTrace Level = "TRACE"
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

// Supported custom slog.Level values.
const (
	SlogLevelTrace slog.Level = -8
	SlogLevelDebug slog.Level = slog.LevelDebug // -4
	SlogLevelInfo  slog.Level = slog.LevelInfo  // 0
	SlogLevelWarn  slog.Level = slog.LevelWarn  // 4
	SlogLevelError slog.Level = slog.LevelError // 8
	SlogLevelFatal slog.Level = 12
)

// String returns the string representation of the Level.
func (l Level) String() string {
	return string(l)
}

// SlogLevel converts Level to its corresponding slog.Level value.
func (l Level) SlogLevel() (slog.Level, error) {
	return ParseSlogLevel(string(l))
}

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
