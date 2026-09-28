package logging

import (
	"fmt"
	"log/slog"
	"strings"
)

// matchTraceDebug parses TRACE or DEBUG level.
func matchTraceDebug(norm string) (Level, slog.Level, bool) {
	if norm == "trace" {
		return LevelTrace, SlogLevelTrace, true
	}
	if norm == "debug" {
		return LevelDebug, SlogLevelDebug, true
	}
	return "", 0, false
}

// isWarn checks if normalized string is a warn alias.
func isWarn(norm string) bool {
	switch norm {
	case "warn", "warning":
		return true
	}
	return false
}

// matchInfoWarn parses INFO or WARN level.
func matchInfoWarn(norm string) (Level, slog.Level, bool) {
	if norm == "info" {
		return LevelInfo, SlogLevelInfo, true
	}
	if isWarn(norm) {
		return LevelWarn, SlogLevelWarn, true
	}
	return "", 0, false
}

// isFatalString checks if norm string is a fatal alias.
func isFatalString(norm string) bool {
	switch norm {
	case "fatal", "critical", "crit", "panic":
		return true
	}
	return false
}

// matchErrorFatal parses ERROR or FATAL level.
func matchErrorFatal(norm string) (Level, slog.Level, bool) {
	if norm == "error" {
		return LevelError, SlogLevelError, true
	}
	if isFatalString(norm) {
		return LevelFatal, SlogLevelFatal, true
	}
	return "", 0, false
}

// matchLevel maps a normalized string to Level and slog.Level.
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
