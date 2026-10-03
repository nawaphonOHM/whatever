package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSlogLevel_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"trace", "trace", SlogLevelTrace},
		{"TRACE", testTrace, SlogLevelTrace},
		{"debug", "debug", SlogLevelDebug},
		{"DEBUG", testDebug, SlogLevelDebug},
		{"info", "info", SlogLevelInfo},
		{"INFO", testInfo, SlogLevelInfo},
		{"empty", "", SlogLevelInfo},
		{"warn", "warn", SlogLevelWarn},
		{"warning", "warning", SlogLevelWarn},
		{"error", "error", SlogLevelError},
		{"fatal", "fatal", SlogLevelFatal},
		{"critical", "critical", SlogLevelFatal},
		{"crit", "crit", SlogLevelFatal},
		{"panic", "panic", SlogLevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slvl, err := ParseSlogLevel(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, slvl)
		})
	}
}

func TestParseSlogLevel_Invalid(t *testing.T) {
	invalidInputs := []string{
		"unknown",
		"none",
		"all",
	}

	for _, in := range invalidInputs {
		t.Run(in, func(t *testing.T) {
			slvl, err := ParseSlogLevel(in)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidLevel)
			assert.Equal(t, slog.Level(0), slvl)
		})
	}
}

func TestLevel_SlogLevel(t *testing.T) {
	tests := []struct {
		lvl      Level
		expected slog.Level
	}{
		{LevelTrace, SlogLevelTrace},
		{LevelDebug, SlogLevelDebug},
		{LevelInfo, SlogLevelInfo},
		{LevelWarn, SlogLevelWarn},
		{LevelError, SlogLevelError},
		{LevelFatal, SlogLevelFatal},
	}

	for _, tt := range tests {
		t.Run(string(tt.lvl), func(t *testing.T) {
			slvl, err := tt.lvl.SlogLevel()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, slvl)
		})
	}

	t.Run("invalid level SlogLevel", func(t *testing.T) {
		invalid := Level("INVALID")
		_, err := invalid.SlogLevel()
		require.ErrorIs(t, err, ErrInvalidLevel)
	})
}
