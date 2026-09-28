package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type parseLevelCase struct {
	expected Level
	input    string
}

func parseLevelValidCases() []parseLevelCase {
	return []parseLevelCase{
		{expected: LevelTrace, input: "trace"},
		{expected: LevelTrace, input: "TRACE"},
		{expected: LevelTrace, input: "Trace"},
		{expected: LevelTrace, input: "  trace  "},
		{expected: LevelDebug, input: "debug"},
		{expected: LevelDebug, input: "DEBUG"},
		{expected: LevelInfo, input: "info"},
		{expected: LevelInfo, input: "INFO"},
		{expected: LevelInfo, input: ""},
		{expected: LevelInfo, input: "   "},
		{expected: LevelWarn, input: "warn"},
		{expected: LevelWarn, input: "WARN"},
		{expected: LevelWarn, input: "warning"},
		{expected: LevelError, input: "error"},
		{expected: LevelError, input: "ERROR"},
		{expected: LevelFatal, input: "fatal"},
		{expected: LevelFatal, input: "FATAL"},
		{expected: LevelFatal, input: "critical"},
		{expected: LevelFatal, input: "crit"},
		{expected: LevelFatal, input: "panic"},
	}
}

func TestParseLevel_Valid(t *testing.T) {
	for _, tt := range parseLevelValidCases() {
		t.Run(tt.input, func(t *testing.T) {
			lvl, err := ParseLevel(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, lvl)
		})
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	invalid := []string{"invalid", "unknown", "verbose", "notice", "123"}
	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := ParseLevel(in)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidLevel)
		})
	}
}
