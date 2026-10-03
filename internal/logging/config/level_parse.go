package config

import (
	"fmt"
	"log/slog"
	"strings"
)

func matchTraceDebug(norm string) (Level, slog.Level, bool) {
	if norm == "trace" {
		return LevelTrace, SlogLevelTrace, true
	}
	if norm == "debug" {
		return LevelDebug, SlogLevelDebug, true
	}
	return "", 0, false
}

func isWarn(norm string) bool {
	return norm == "warn" || norm == "warning"
}

func matchInfoWarn(norm string) (Level, slog.Level, bool) {
	if norm == "info" {
		return LevelInfo, SlogLevelInfo, true
	}
	if isWarn(norm) {
		return LevelWarn, SlogLevelWarn, true
	}
	return "", 0, false
}

func isFatalString(norm string) bool {
	switch norm {
	case "fatal", "critical", "crit", "panic":
		return true
	}
	return false
}

func matchErrorFatal(norm string) (Level, slog.Level, bool) {
	if norm == "error" {
		return LevelError, SlogLevelError, true
	}
	if isFatalString(norm) {
		return LevelFatal, SlogLevelFatal, true
	}
	return "", 0, false
}

func matchLevel(norm string) (Level, slog.Level, bool) {
	if lvl, slvl, ok := matchTraceDebug(norm); ok {
		return lvl, slvl, true
	}
	if lvl, slvl, ok := matchInfoWarn(norm); ok {
		return lvl, slvl, true
	}
	return matchErrorFatal(norm)
}

// ParseLevel parses a case-insensitive string into a Level.
func ParseLevel(s string) (Level, error) {
	trimmed := strings.ToLower(strings.TrimSpace(s))
	if trimmed == "" {
		return LevelInfo, nil
	}
	if lvl, _, ok := matchLevel(trimmed); ok {
		return lvl, nil
	}
	return "", fmt.Errorf("%w: %q (supported: TRACE, DEBUG, INFO, WARN, ERROR, FATAL)", ErrInvalidLevel, s)
}

// ParseSlogLevel parses a case-insensitive string into a slog.Level.
func ParseSlogLevel(s string) (slog.Level, error) {
	trimmed := strings.ToLower(strings.TrimSpace(s))
	if trimmed == "" {
		return SlogLevelInfo, nil
	}
	if _, slvl, ok := matchLevel(trimmed); ok {
		return slvl, nil
	}
	return 0, fmt.Errorf("%w: %q (supported: TRACE, DEBUG, INFO, WARN, ERROR, FATAL)", ErrInvalidLevel, s)
}
