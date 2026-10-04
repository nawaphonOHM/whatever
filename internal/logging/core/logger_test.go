package core

import (
	"bytes"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Defaults(t *testing.T) {
	l := New()
	require.NotNil(t, l)
	assert.NotNil(t, l.Slog())
	assert.NotNil(t, l.Handler())
	assert.NotNil(t, l.ExitFunc())

	fallbackL := New(&config.Config{Output: "invalid://err"})
	require.NotNil(t, fallbackL)
	assert.NotNil(t, fallbackL.Slog())
}

func TestNewFromConfig(t *testing.T) {
	t.Run("nil config uses defaults", func(t *testing.T) {
		l, err := NewFromConfig(nil)
		require.NoError(t, err)
		require.NotNil(t, l)
	})

	t.Run("invalid output errors", func(t *testing.T) {
		cfg := &config.Config{Output: "invalid://output"}
		l, err := NewFromConfig(cfg)
		assert.Error(t, err)
		assert.Nil(t, l)
	})
}

func TestNew_LevelFiltering(t *testing.T) {
	for _, tt := range levelFilterCases() {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := NewJSON(buf, tt.cfgLevel)
			tt.logAction(l)

			if !tt.expectLog {
				assert.Empty(t, buf.String())
				return
			}
			logMap := parseJSONLog(t, buf.Bytes())
			assert.Equal(t, tt.expectedLevel, logMap["level"])
		})
	}
}

func TestNewJSON_And_NewText(t *testing.T) {
	t.Run("NewJSON", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := NewJSON(buf, config.LevelDebug)
		l.Debug("json debug message")
		logMap := parseJSONLog(t, buf.Bytes())
		assert.Equal(t, "DEBUG", logMap["level"])
		assert.Equal(t, "json debug message", logMap["msg"])
	})

	t.Run("NewText", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := NewText(buf, config.LevelInfo)
		l.Info("text info message")
		assert.Contains(t, buf.String(), "level=INFO")
		assert.Contains(t, buf.String(), "msg=\"text info message\"")
	})

	t.Run("nil writer defaults to stdout", func(t *testing.T) {
		lj := NewJSON(nil, config.LevelInfo)
		require.NotNil(t, lj)
		lt := NewText(nil, config.LevelInfo)
		require.NotNil(t, lt)
	})
}
