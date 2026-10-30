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
	testHTTPStatus204Num = 204
	testRedactedVal      = "[REDACTED]"
	testKeyMethodVal     = "GET"
	testPathAPI          = "/api"
)

func redactStatusAttr(a slog.Attr) slog.Attr {
	if a.Key == testKeyStatusName {
		return slog.Int(a.Key, testHTTPStatus204Num)
	}
	return a
}

func redactAgentAttr(a slog.Attr) slog.Attr {
	if a.Key == "user_agent" {
		return slog.Attr{}
	}
	return redactStatusAttr(a)
}

func redactTopAttr(a slog.Attr) slog.Attr {
	if a.Key == "client_ip" {
		return slog.String(a.Key, testRedactedVal)
	}
	return redactAgentAttr(a)
}

func redactHTTPAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	return redactTopAttr(a)
}

func verifyRedactedHTTPLog(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testRedactedVal, httpReq["remoteIp"])
	assert.Empty(t, httpReq["userAgent"])
	assert.Equal(t, float64(testHTTPStatus204Num), httpReq["status"])
	assert.Equal(t, testKeyMethodVal, httpReq[testKeyReqMethod])
	assert.Equal(t, testPathAPI, httpReq[testKeyReqURL])
}

func TestGCPHandler_ReplaceAttrBeforeHTTPExtraction(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		ReplaceAttr: redactHTTPAttr,
	}
	logger := slog.New(core.NewGCPHandler(&buf, opts))
	logger.Info("HTTP Request",
		"method", testKeyMethodVal,
		"path", testPathAPI,
		testKeyStatusName, testHTTPStatus200Num,
		"client_ip", "1.2.3.4",
		"user_agent", "Mozilla/5.0",
	)

	verifyRedactedHTTPLog(t, buf.Bytes())
}
