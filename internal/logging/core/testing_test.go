package core

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetTestLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	testLogger := NewJSON(buf, config.LevelInfo)

	cleanup := SetTestLogger(testLogger)
	assert.Equal(t, testLogger, Default())

	Info("captured test message")
	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "captured test message", logMap["msg"])

	cleanup()
	assert.NotEqual(t, testLogger, Default())
}

func TestCaptureLogs(t *testing.T) {
	t.Run("default level debug", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cleanup := CaptureLogs(buf)
		defer cleanup()

		Debug("debug message")
		Info("info message")

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		require.Len(t, lines, 2)
		assert.Contains(t, lines[0], "debug message")
		assert.Contains(t, lines[1], "info message")
	})

	t.Run("custom level warn", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cleanup := CaptureLogs(buf, config.LevelWarn)
		defer cleanup()

		Info("suppressed info message")
		Warn("warn message")

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		require.Len(t, lines, 1)
		assert.Contains(t, lines[0], "warn message")
	})
}

func TestCaptureTextLogs(t *testing.T) {
	t.Run("default level debug", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cleanup := CaptureTextLogs(buf)
		defer cleanup()

		Debug("debug text message")
		assert.Contains(t, buf.String(), "level=DEBUG")
		assert.Contains(t, buf.String(), "msg=\"debug text message\"")
	})

	t.Run("custom level warn", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cleanup := CaptureTextLogs(buf, config.LevelWarn)
		defer cleanup()

		Info("suppressed info text message")
		Warn("warn text message")

		assert.NotContains(t, buf.String(), "suppressed info text message")
		assert.Contains(t, buf.String(), "level=WARN")
		assert.Contains(t, buf.String(), "msg=\"warn text message\"")
	})
}

func TestResetDefault(t *testing.T) {
	buf := &bytes.Buffer{}
	testLogger := NewJSON(buf, config.LevelInfo)
	SetDefault(testLogger)
	assert.Equal(t, testLogger, Default())

	ResetDefault()
	assert.NotEqual(t, testLogger, Default())
}

func TestResetDefault_WithFileCloser(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "reset_app.log")

	cfg := &config.Config{
		Output: logPath,
		Level:  "info",
		Format: "text",
	}

	l, err := NewFromConfig(cfg)
	require.NoError(t, err)
	SetDefault(l)

	Info("message before reset")
	ResetDefault()
	assert.NotEqual(t, l, Default())
}
