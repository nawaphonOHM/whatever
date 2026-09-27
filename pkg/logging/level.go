// Package logging provides general-purpose structured logging on top of log/slog
// with OpenTelemetry distributed trace correlation, configurable log levels
// (TRACE, DEBUG, INFO, WARN, ERROR, FATAL), and customizable output formats.
package logging

import (
	"errors"
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

// ErrInvalidLevel indicates an unrecognized log level string.
var ErrInvalidLevel = errors.New("invalid log level")

// String returns the string representation of the Level.
func (l Level) String() string {
	return string(l)
}

// SlogLevel converts Level to its corresponding slog.Level value.
func (l Level) SlogLevel() (slog.Level, error) {
	return ParseSlogLevel(string(l))
}
