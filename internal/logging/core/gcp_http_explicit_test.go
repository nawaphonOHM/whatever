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
	testHTTPStatus200Num = 200
	testMethodGET        = "GET"
	testKeyMethod        = "method"
	testKeyStatus        = "status"
	testKeyFoo           = "foo"
)

func verifyExplicitGroupLeftover(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGET, httpReq[testKeyReqMethod])

	customReq, ok := entry["custom_httpRequest"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "done", customReq[testKeyStatus])
	assert.Equal(t, "bar", customReq[testKeyFoo])
}

func TestGCPHandler_ExplicitGroup_LeftoverPreserved(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	group := slog.Group("httpRequest", testKeyMethod, testMethodGET, testKeyStatus, "done", testKeyFoo, "bar")
	logger.Info("processed", group)

	verifyExplicitGroupLeftover(t, buf.Bytes())
}

func verifyExplicitMapLeftover(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGET, httpReq[testKeyReqMethod])
	assert.Equal(t, "0.015s", httpReq[testKeyReqLatency])
	assert.Equal(t, float64(testHTTPStatus200Num), httpReq[testKeyStatus])

	customReq, ok := entry["custom_httpRequest"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(1), customReq[testKeyFoo])
}

func TestGCPHandler_ExplicitMap_LeftoverAndNormalization(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	m := map[string]any{
		testKeyMethod: testMethodGET,
		"latency":     "15ms",
		testKeyStatus: "200",
		testKeyFoo:    1,
	}
	logger.Info("processed", "httpRequest", m)

	verifyExplicitMapLeftover(t, buf.Bytes())
}

func TestGCPHandler_ExplicitStruct_LatencyNormalization(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	explicitReq := &core.HttpRequest{
		RequestMethod: testMethodGET,
		Latency:       "15ms",
	}
	logger.Info("processed", "httpRequest", explicitReq)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "0.015s", httpReq[testKeyReqLatency])
}
