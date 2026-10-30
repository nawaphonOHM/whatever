package core_test

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

type stringValuer struct{ val string }

func (v stringValuer) LogValue() slog.Value { return slog.StringValue(v.val) }

type statusValuer struct{ st int }

func (v statusValuer) LogValue() slog.Value { return slog.IntValue(v.st) }

type latencyValuer struct{ d time.Duration }

func (v latencyValuer) LogValue() slog.Value { return slog.DurationValue(v.d) }

const (
	testRedactedStr      = "[REDACTED]"
	testStatus200Num     = 200
	testLatencyFifteen   = 15 * time.Millisecond
	testStatus99999      = 99999
	testStatus1Num       = 1
	testStatusFractional = 200.5
	testMethodInt42      = 42
	testKeyReqRemoteIP   = "remoteIp"
	testKeyReqUserAgent  = "userAgent"
)

func verifyLogValuerExtraction(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testRedactedStr, httpReq[testKeyReqRemoteIP])
	assert.Equal(t, float64(testStatus200Num), httpReq[testKeyStatusStr])
	assert.Equal(t, "0.015s", httpReq[testKeyReqLatency])
}

func TestGCPHandler_LogValuer_ResolvedBeforeConversion(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", testKeyMethodStr, testMethodGETStr, testKeyPathStr, "/",
		"client_ip", stringValuer{val: testRedactedStr},
		testKeyStatusStr, statusValuer{st: testStatus200Num},
		"latency", latencyValuer{d: testLatencyFifteen})

	verifyLogValuerExtraction(t, buf.Bytes())
}

func verifyStatusPayloadMatch(t *testing.T, expected, actual any) {
	if n, ok := expected.(int); ok {
		assert.Equal(t, float64(n), actual)
		return
	}
	assert.Equal(t, expected, actual)
}

func verifyInvalidStatusPreserved(t *testing.T, rawStatus any) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", testKeyMethodStr, testMethodGETStr, testKeyPathStr, "/", testKeyStatusStr, rawStatus)

	entry := parseLogMap(t, buf.Bytes())
	verifyStatusPayloadMatch(t, rawStatus, entry[testKeyStatusStr])
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Empty(t, httpReq[testKeyStatusStr])
}

func TestGCPHandler_InvalidAndOutOfRangeStatus_PreservedInPayload(t *testing.T) {
	verifyInvalidStatusPreserved(t, testStatus99999)
	verifyInvalidStatusPreserved(t, testStatus1Num)
	verifyInvalidStatusPreserved(t, testStatusFractional)
	verifyInvalidStatusPreserved(t, "done")
}

func verifyNonStringPayload(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	assert.Equal(t, float64(testMethodInt42), entry[testKeyMethodStr])
	assert.Equal(t, []any{"custom", "agent"}, entry["user_agent"])
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Empty(t, httpReq[testKeyReqMethod])
	assert.Empty(t, httpReq[testKeyReqUserAgent])
}

func TestGCPHandler_NonStringMethodAndAgent_PreservedInPayload(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	uaSlice := []string{"custom", "agent"}
	logger.Info("HTTP Request", testKeyPathStr, "/api", testKeyMethodStr, testMethodInt42, "user_agent", uaSlice)
	verifyNonStringPayload(t, buf.Bytes())
}
