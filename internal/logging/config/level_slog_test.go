package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testSlogBelowTrace             slog.Level = -10
	testSlogBetweenTraceDebug      slog.Level = -6
	testSlogBetweenDebugInfo       slog.Level = -2
	testSlogBetweenInfoWarn        slog.Level = 2
	testSlogLevelBetweenWarnError  slog.Level = 6
	testSlogLevelBetweenErrorFatal slog.Level = 10
	testSlogLevelAboveFatal        slog.Level = 14
)

func TestLevelFromSlog(t *testing.T) {
	tests := []struct {
		name     string
		expected Level
		input    slog.Level
	}{
		{"below trace", LevelTrace, testSlogBelowTrace},
		{"trace level", LevelTrace, SlogLevelTrace},
		{"between trace and debug", LevelDebug, testSlogBetweenTraceDebug},
		{"debug level", LevelDebug, SlogLevelDebug},
		{"between debug and info", LevelDebug, testSlogBetweenDebugInfo},
		{"info level", LevelInfo, SlogLevelInfo},
		{"between info and warn", LevelInfo, testSlogBetweenInfoWarn},
		{"warn level", LevelWarn, SlogLevelWarn},
		{"between warn and error", LevelWarn, testSlogLevelBetweenWarnError},
		{"error level", LevelError, SlogLevelError},
		{"between error and fatal", LevelError, testSlogLevelBetweenErrorFatal},
		{"fatal level", LevelFatal, SlogLevelFatal},
		{"above fatal", LevelFatal, testSlogLevelAboveFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl := LevelFromSlog(tt.input)
			assert.Equal(t, tt.expected, lvl)
		})
	}
}
