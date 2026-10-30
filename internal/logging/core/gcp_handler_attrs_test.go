package core_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testUserSeverity  = "USER_SEVERITY"
	testUserTimestamp = "USER_TIMESTAMP"
	testUserMessage   = "USER_MESSAGE"
	testUserHTTP      = "USER_HTTP"
	testUserID101     = 101
)

func TestGCPHandlerReservedKeyCollisions(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler)

	logger.Info("reserved test",
		"severity", testUserSeverity,
		"timestamp", testUserTimestamp,
		"message", testUserMessage,
		"httpRequest", testUserHTTP,
	)

	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, testUserSeverity, entry["custom_severity"])
	assert.Equal(t, testUserTimestamp, entry["custom_timestamp"])
	assert.Equal(t, testUserMessage, entry["custom_message"])
	assert.Equal(t, testUserHTTP, entry["custom_httpRequest"])
}

func TestGCPHandlerWithAttrsAndGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler).With("service", "auth").WithGroup("audit")

	logger.Info("user action", "user_id", testUserID101)

	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "auth", entry["service"])
	auditGroup, ok := entry["audit"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(testUserID101), auditGroup["user_id"])
}

func TestGCPHandlerLevelEnabled(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})

	assert.False(t, handler.Enabled(context.Background(), slog.LevelDebug))
	assert.True(t, handler.Enabled(context.Background(), slog.LevelWarn))
	assert.True(t, handler.Enabled(context.Background(), slog.LevelError))
}
