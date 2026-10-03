package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLevel_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
	}{
		{"trace lower", "trace", LevelTrace},
		{"trace upper", testTrace, LevelTrace},
		{"trace spaces", "  trace  ", LevelTrace},
		{"debug lower", "debug", LevelDebug},
		{"debug upper", testDebug, LevelDebug},
		{"info lower", "info", LevelInfo},
		{"info upper", testInfo, LevelInfo},
		{"info empty", "", LevelInfo},
		{"info whitespace", "   ", LevelInfo},
		{"warn lower", "warn", LevelWarn},
		{"warn upper", testWarn, LevelWarn},
		{"warning lower", "warning", LevelWarn},
		{"warning upper", "WARNING", LevelWarn},
		{"error lower", "error", LevelError},
		{"error upper", testError, LevelError},
		{"fatal lower", "fatal", LevelFatal},
		{"fatal upper", testFatal, LevelFatal},
		{"critical lower", "critical", LevelFatal},
		{"critical upper", "CRITICAL", LevelFatal},
		{"crit lower", "crit", LevelFatal},
		{"crit upper", "CRIT", LevelFatal},
		{"panic lower", "panic", LevelFatal},
		{"panic upper", "PANIC", LevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl, err := ParseLevel(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, lvl)
		})
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	invalidInputs := []string{
		"unknown",
		"verbose",
		"err",
		"inf",
		"123",
	}

	for _, in := range invalidInputs {
		t.Run(in, func(t *testing.T) {
			lvl, err := ParseLevel(in)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidLevel)
			assert.Empty(t, lvl)
		})
	}
}

func TestLevel_String(t *testing.T) {
	assert.Equal(t, testTrace, LevelTrace.String())
	assert.Equal(t, testDebug, LevelDebug.String())
	assert.Equal(t, testInfo, LevelInfo.String())
	assert.Equal(t, testWarn, LevelWarn.String())
	assert.Equal(t, testError, LevelError.String())
	assert.Equal(t, testFatal, LevelFatal.String())
}
