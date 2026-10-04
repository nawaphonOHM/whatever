package core

import (
	"io"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
)

// SetTestLogger sets the default logger to the provided logger
// and returns a cleanup function to restore the previous logger.
func SetTestLogger(l *Logger) func() {
	prev := Default()
	SetDefault(l)
	return func() {
		if l != nil {
			l.Flush()
		}
		SetDefault(prev)
	}
}

// CaptureLogs redirects default logging to the provided writer in JSON format
// at the specified level and returns a cleanup function to restore the previous logger.
func CaptureLogs(w io.Writer, lvl ...config.Level) func() {
	level := config.LevelDebug
	if len(lvl) > 0 && lvl[0] != "" {
		level = lvl[0]
	}
	testLogger := NewJSON(w, level)
	cleanup := SetTestLogger(testLogger)
	return func() {
		testLogger.Flush()
		cleanup()
	}
}

// CaptureTextLogs redirects default logging to the provided writer in Text format
// at the specified level and returns a cleanup function to restore the previous logger.
func CaptureTextLogs(w io.Writer, lvl ...config.Level) func() {
	level := config.LevelDebug
	if len(lvl) > 0 && lvl[0] != "" {
		level = lvl[0]
	}
	testLogger := NewText(w, level)
	cleanup := SetTestLogger(testLogger)
	return func() {
		testLogger.Flush()
		cleanup()
	}
}
