package core_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testKeyMessage     = "message"
	testKeySeverity    = "severity"
	testKeyTimestamp   = "timestamp"
	testKeyHTTPRequest = "httpRequest"
	testKeyReqMethod   = "requestMethod"
	testKeyReqURL      = "requestUrl"
	testKeyReqLatency  = "latency"
	testSampleLatency  = 15 * time.Millisecond
	testBytes1024      = 1024
	testStatus200      = 200
)

func parseLogMap(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var entry map[string]any
	err := json.Unmarshal(data, &entry)
	require.NoError(t, err)
	return entry
}

func TestGCPHandlerBasic(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	logger.Info("service started", "env", "prod", "version", 1)

	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, core.GCPSeverityInfo, entry[testKeySeverity])
	assert.Equal(t, "service started", entry[testKeyMessage])
	assert.Equal(t, "prod", entry["env"])
	assert.NotNil(t, entry[testKeyTimestamp])
}

func TestGCPHandlerJSONMessage(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler)

	logger.Info(`{"action": "create_user", "id": 42}`)

	entry := parseLogMap(t, buf.Bytes())
	msgObj, ok := entry[testKeyMessage].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "create_user", msgObj["action"])
}

func TestGCPHandlerHTTPLogging(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler)

	logger.Info("HTTP Request",
		"method", "GET",
		"path", "/api/v1/users",
		"query", "limit=10",
		"status", testStatus200,
		"latency", testSampleLatency,
		"client_ip", "192.168.1.1",
		"user_agent", "curl/7.68.0",
		"bytes_out", testBytes1024,
	)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "GET", httpReq[testKeyReqMethod])
	assert.Equal(t, "/api/v1/users?limit=10", httpReq[testKeyReqURL])
	assert.Equal(t, "0.015s", httpReq[testKeyReqLatency])
}
