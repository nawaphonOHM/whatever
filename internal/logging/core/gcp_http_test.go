package core_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testHTTPStatus201 = 201
	testHTTPLatencyMs = 125
)

func verifyExplicitHTTPLog(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "POST", httpReq[testKeyReqMethod])
	assert.Equal(t, "/submit", httpReq[testKeyReqURL])
	assert.Equal(t, "value", entry["extra_attr"])
}

func TestGCPHandlerHTTPExplicitStruct(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	reqObj := &core.HttpRequest{
		RequestMethod: "POST",
		RequestURL:    "/submit",
		Status:        testHTTPStatus201,
		Latency:       "0.050s",
	}
	logger.Info("processed", testKeyHTTPRequest, reqObj, "extra_attr", "value")
	verifyExplicitHTTPLog(t, buf.Bytes())
}

func TestGCPHandlerUnrelatedStatusNotHTTP(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("database health check", "status", "healthy", "service", "db")

	entry := parseLogMap(t, buf.Bytes())
	assert.Nil(t, entry[testKeyHTTPRequest])
	assert.Equal(t, "healthy", entry["status"])
}

func TestGCPHandlerNumericLatency(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", "method", "GET", "path", "/items", "latency_ms", testHTTPLatencyMs)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "0.125s", httpReq[testKeyReqLatency])
}
