package logging

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatLevelAttr(t *testing.T) {
	tests := []struct {
		input    slog.Attr
		expected slog.Attr
		name     string
	}{
		{
			name:     "slog LevelInfo converted to INFO string",
			input:    slog.Any(slog.LevelKey, slog.LevelInfo),
			expected: slog.String(slog.LevelKey, "INFO"),
		},
		{
			name:     "slog LevelTrace converted to TRACE string",
			input:    slog.Any(slog.LevelKey, SlogLevelTrace),
			expected: slog.String(slog.LevelKey, "TRACE"),
		},
		{
			name:     "slog LevelFatal converted to FATAL string",
			input:    slog.Any(slog.LevelKey, SlogLevelFatal),
			expected: slog.String(slog.LevelKey, "FATAL"),
		},
		{
			name:     "non-level value passes through unchanged",
			input:    slog.String(slog.LevelKey, "already-string"),
			expected: slog.String(slog.LevelKey, "already-string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := formatLevelAttr(tt.input)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestBuildReplaceAttr_NilCustom(t *testing.T) {
	replacer := buildReplaceAttr(nil)

	levelAttr := replacer(nil, slog.Any(slog.LevelKey, slog.LevelWarn))
	assert.Equal(t, slog.String(slog.LevelKey, "WARN"), levelAttr)

	msgAttr := replacer(nil, slog.String("msg", "hello world"))
	assert.Equal(t, slog.String("msg", "hello world"), msgAttr)

	flagAttr := replacer([]string{"group1"}, slog.Bool("active", true))
	assert.Equal(t, slog.Bool("active", true), flagAttr)
}

func TestBuildReplaceAttr_CustomChained(t *testing.T) {
	custom := func(groups []string, a slog.Attr) slog.Attr {
		assert.Equal(t, []string{"grp"}, groups)
		if a.Key == slog.LevelKey {
			return slog.String(a.Key, "PREFIX_"+a.Value.String())
		}
		return a
	}

	replacer := buildReplaceAttr(custom)
	levelRes := replacer([]string{"grp"}, slog.Any(slog.LevelKey, slog.LevelInfo))
	assert.Equal(t, slog.String(slog.LevelKey, "PREFIX_INFO"), levelRes)
}

func TestBuildReplaceAttr_CustomChainedNonLevel(t *testing.T) {
	custom := func(groups []string, a slog.Attr) slog.Attr {
		assert.Empty(t, groups)
		if a.Key == "password" {
			return slog.String(a.Key, "[REDACTED]")
		}
		return a
	}

	replacer := buildReplaceAttr(custom)
	secretRes := replacer(nil, slog.String("password", "supersecret"))
	assert.Equal(t, slog.String("password", "[REDACTED]"), secretRes)
}
