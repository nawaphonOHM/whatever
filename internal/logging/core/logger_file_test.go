package core

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyParsedLogPair(t *testing.T, line1, line2 string) {
	t.Helper()
	log1 := parseJSONLog(t, []byte(line1))
	assert.Equal(t, "log message", log1[testMsg])
	assert.Equal(t, "v1", log1["k1"])
	assert.Equal(t, "INFO", log1["level"])

	log2 := parseJSONLog(t, []byte(line2))
	assert.Equal(t, "log attrs message", log2[testMsg])
	assert.Equal(t, "v2", log2["k2"])
	assert.Equal(t, "WARN", log2["level"])
}

func TestLogger_LogAndLogAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	l := NewJSON(buf, config.LevelDebug)

	l.Log(context.Background(), config.SlogLevelInfo, "log message", "k1", "v1")
	l.LogAttrs(context.Background(), config.SlogLevelWarn, "log attrs message", slog.String("k2", "v2"))

	lines := strings.Split(strings.TrimSpace(buf.String()), testNewline)
	require.Len(t, lines, 2)
	verifyParsedLogPair(t, lines[0], lines[1])
}

func TestLogger_FileOutputAndClose(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "app.log")
	cfg := &config.Config{Output: logPath, Level: "info", Format: "text"}

	l, err := NewFromConfig(cfg)
	require.NoError(t, err)
	require.NotNil(t, l)

	l.Info("file message 1")
	require.NoError(t, l.Close())

	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "msg=\"file message 1\"")
}

func verifyGlobalParsedLogs(t *testing.T, line1, line2 string) {
	t.Helper()
	log1 := parseJSONLog(t, []byte(line1))
	assert.Equal(t, "pkg log msg", log1[testMsg])
	assert.Equal(t, "bar", log1["foo"])

	log2 := parseJSONLog(t, []byte(line2))
	assert.Equal(t, "pkg log attrs msg", log2[testMsg])
	assert.Equal(t, "dog", log2["cat"])
}

func TestGlobal_LogAndLogAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	cleanup := CaptureLogs(buf, config.LevelDebug)
	defer cleanup()

	Log(context.Background(), config.SlogLevelInfo, "pkg log msg", "foo", "bar")
	LogAttrs(context.Background(), config.SlogLevelWarn, "pkg log attrs msg", slog.String("cat", "dog"))

	lines := strings.Split(strings.TrimSpace(buf.String()), testNewline)
	require.Len(t, lines, 2)
	verifyGlobalParsedLogs(t, lines[0], lines[1])
}
