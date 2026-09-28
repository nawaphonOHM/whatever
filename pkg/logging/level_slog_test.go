package logging

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type parseSlogCase struct {
	input    string
	expected slog.Level
}

func parseSlogCases() []parseSlogCase {
	return []parseSlogCase{
		{input: "trace", expected: SlogLevelTrace},
		{input: "TRACE", expected: SlogLevelTrace},
		{input: "debug", expected: SlogLevelDebug},
		{input: "info", expected: SlogLevelInfo},
		{input: "", expected: SlogLevelInfo},
		{input: "warn", expected: SlogLevelWarn},
		{input: "error", expected: SlogLevelError},
		{input: "fatal", expected: SlogLevelFatal},
	}
}

func TestParseSlogLevel_Valid(t *testing.T) {
	for _, tt := range parseSlogCases() {
		t.Run(tt.input, func(t *testing.T) {
			lvl, err := ParseSlogLevel(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, lvl)
		})
	}
}

func TestParseSlogLevel_Invalid(t *testing.T) {
	_, err := ParseSlogLevel("invalid")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidLevel)
}

func TestLevel_StringAndSlogLevel(t *testing.T) {
	assert.Equal(t, "TRACE", LevelTrace.String())
	assert.Equal(t, "FATAL", LevelFatal.String())

	sl, err := LevelTrace.SlogLevel()
	require.NoError(t, err)
	assert.Equal(t, SlogLevelTrace, sl)

	invalid := Level("unknown")
	_, err = invalid.SlogLevel()
	require.Error(t, err)
}

type slogToLevelCase struct {
	expected Level
	input    slog.Level
}

func slogToLevelCases() []slogToLevelCase {
	return []slogToLevelCase{
		{expected: LevelTrace, input: SlogLevelTrace},
		{expected: LevelDebug, input: SlogLevelDebug},
		{expected: LevelInfo, input: SlogLevelInfo},
		{expected: LevelWarn, input: SlogLevelWarn},
		{expected: LevelError, input: SlogLevelError},
		{expected: LevelFatal, input: SlogLevelFatal},
	}
}

func TestLevelFromSlog(t *testing.T) {
	for _, tt := range slogToLevelCases() {
		t.Run(tt.input.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, LevelFromSlog(tt.input))
		})
	}
}
