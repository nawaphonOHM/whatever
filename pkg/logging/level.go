// Package logging provides general-purpose structured logging functions
// driven by environment variable configuration.
package logging

import (
	"github.com/nawaphonOHM/whatever/internal/logging/config"
)

// Level represents supported logging severity levels.
type Level = config.Level

// Supported log level string constants.
const (
	LevelTrace Level = config.LevelTrace
	LevelDebug Level = config.LevelDebug
	LevelInfo  Level = config.LevelInfo
	LevelWarn  Level = config.LevelWarn
	LevelError Level = config.LevelError
	LevelFatal Level = config.LevelFatal
)
