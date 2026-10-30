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

const (
	testLatencyPointOne      = "0.1s"
	testLatencyPointZeroFive = "0.05s"
	testLatencyFive          = 5 * time.Millisecond
	testLatency15Num         = 15
	testLatencyFifty         = 50 * time.Millisecond
	testQueryA1              = "a=1"
	testQueryB2              = "b=2"
	testKeyQuery             = "query"
	testKeyLatency           = "latency"
	testKeyLatencyMs         = "latency_ms"
)

func verifyExplicitLatencyPreserved(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testLatencyPointOne, httpReq[testKeyReqLatency])
	assert.Equal(t, "5ms", entry[testKeyLatency])
}

func TestGCPHandler_ExplicitLatency_DoesNotOverwrite(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", testKeyHTTPRequest, core.HttpRequest{Latency: testLatencyPointOne},
		testKeyLatency, testLatencyFive)

	verifyExplicitLatencyPreserved(t, buf.Bytes())
}

func verifyLatencyMsOrdering(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testLatencyPointZeroFive, httpReq[testKeyReqLatency])
	assert.Nil(t, entry[testKeyLatency])
	assert.Nil(t, entry[testKeyLatencyMs])
}

func TestGCPHandler_LatencyMs_Then_Latency_NoLoss(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", testKeyMethodStr, testMethodGETStr, testKeyPathStr, "/",
		testKeyLatencyMs, testLatency15Num, testKeyLatency, testLatencyFifty)

	verifyLatencyMsOrdering(t, buf.Bytes())
}

func verifyMultipleExplicitGroups(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGETStr, httpReq[testKeyReqMethod])

	customReq, ok := entry["custom_httpRequest"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "a", customReq["foo"])
	assert.Equal(t, "b", customReq["bar"])
	assert.Equal(t, "POST", customReq[testKeyMethodStr])
}

func TestGCPHandler_MultipleExplicitGroups_Merged(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	g1 := slog.Group(testKeyHTTPRequest, testKeyMethodStr, testMethodGETStr, "foo", "a")
	g2 := slog.Group(testKeyHTTPRequest, testKeyMethodStr, "POST", "bar", "b")
	logger.Info("job", g1, g2)

	verifyMultipleExplicitGroups(t, buf.Bytes())
}

func verifyExplicitLeftoverVsFlat(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	assert.Equal(t, testQueryB2, entry[testKeyQuery])
	customReq, ok := entry["custom_httpRequest"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testQueryA1, customReq[testKeyQuery])
}

func TestGCPHandler_ExplicitGroupLeftover_CollidingFlatKey(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("job", slog.Group(testKeyHTTPRequest, testKeyQuery, testQueryA1), testKeyQuery, testQueryB2)

	verifyExplicitLeftoverVsFlat(t, buf.Bytes())
}
