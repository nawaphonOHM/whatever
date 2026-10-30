package core_test

import (
	"bytes"
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testFloatHuge         = 1e20
	testFloatSizeFraction = 12.5
	testStatusFraction    = 200.7
	testStatusLow         = 99
	testStatusHigh        = 999
	testLatencyNeg5       = -5
	testKeyRespSize       = "responseSize"
	testKeyBytesOut       = "bytes_out"
)

func verifyPreservedAttr(t *testing.T, key string, val any, targetKey string) {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", "method", "GET", "path", "/", key, val)

	entry := parseLogMapNumbers(t, buf.Bytes())
	require.NotNil(t, entry[key])
	matchPreservedValue(t, val, entry[key])
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	assert.True(t, ok)
	assert.Empty(t, httpReq[targetKey])
}

func TestGCPHandler_LatencyOverflow_PreservedInPayload(t *testing.T) {
	verifyPreservedAttr(t, testKeyLatencyMs, int64(math.MaxInt64), testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, uint64(math.MaxUint64), testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, testFloatHuge, testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, testLatencyNeg5*time.Millisecond, testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, "-5ms", testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, "NaN", testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, "fast", testKeyReqLatency)
	verifyPreservedAttr(t, testKeyLatency, testLatencyNeg5, testKeyReqLatency)
}

func TestGCPHandler_ResponseSizeOverflow_PreservedInPayload(t *testing.T) {
	verifyPreservedAttr(t, testKeyBytesOut, uint64(math.MaxUint64), testKeyRespSize)
	verifyPreservedAttr(t, testKeyBytesOut, testFloatHuge, testKeyRespSize)
	verifyPreservedAttr(t, testKeyBytesOut, testFloatSizeFraction, testKeyRespSize)
	verifyPreservedAttr(t, testKeyBytesOut, "-3", testKeyRespSize)
	verifyPreservedAttr(t, testKeyRespSize, -1, testKeyRespSize)
	verifyPreservedAttr(t, testKeyBytesOut, -1, testKeyRespSize)
}

func TestGCPHandler_StatusInvalid_PreservedInPayload(t *testing.T) {
	verifyPreservedAttr(t, testKeyStatusName, uint64(math.MaxUint64), testKeyStatusName)
	verifyPreservedAttr(t, testKeyStatusName, testStatusFraction, testKeyStatusName)
	verifyPreservedAttr(t, testKeyStatusName, testStatusLow, testKeyStatusName)
	verifyPreservedAttr(t, testKeyStatusName, testStatusHigh, testKeyStatusName)
}
