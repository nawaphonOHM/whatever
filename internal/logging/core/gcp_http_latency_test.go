package core_test

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testFloatLatency15_5 = 15.5
	testFloatLatency0_33 = 0.3333
	testLatency15MsNum   = 15
	testLatency15MsStr   = "0.015s"
	testLatency50MsNum   = 50
	testNanosAlmostSec   = 999999999
	testKeyLatencyCamel  = "latencyMs"
)

func verifyLatencyAttr(t *testing.T, key string, rawVal any, expectedLatency string) {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", "method", "GET", "path", "/", key, rawVal)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, expectedLatency, httpReq[testKeyReqLatency])
}

func verifyLatency(t *testing.T, rawVal any, expectedLatency string) {
	t.Helper()
	verifyLatencyAttr(t, testKeyLatency, rawVal, expectedLatency)
}

func TestGCPHandler_LatencyFormats(t *testing.T) {
	verifyLatency(t, "15ms", testLatency15MsStr)
	verifyLatency(t, "15s", "15s")
	verifyLatency(t, "2m3s", "123s")
	verifyLatency(t, "1.5µs", "0.0000015s")
	verifyLatency(t, "0.045s", "0.045s")
	verifyLatency(t, testFloatLatency15_5, "0.0155s")
	verifyLatency(t, testFloatLatency0_33, "0.0003333s")
	verifyLatency(t, "0.3333", "0.0003333s")
	verifyLatency(t, testLatency15MsNum, testLatency15MsStr)
	verifyLatency(t, 0, "0s")
}

func TestGCPHandler_LatencyMs_KeyAwareFormats(t *testing.T) {
	verifyLatencyAttr(t, testKeyLatencyMs, testFloatLatency15_5, "0.0155s")
	verifyLatencyAttr(t, testKeyLatencyMs, testFloatLatency0_33, "0.0003333s")
	verifyLatencyAttr(t, testKeyLatencyMs, "15", testLatency15MsStr)
	verifyLatencyAttr(t, testKeyLatencyMs, "15.5", "0.0155s")
	verifyLatencyAttr(t, testKeyLatencyMs, "15ms", testLatency15MsStr)
	verifyLatencyAttr(t, testKeyLatencyCamel, testLatency50MsNum, "0.05s")
}

func TestGCPHandler_LatencyDurationPrecision(t *testing.T) {
	verifyLatency(t, 1*time.Nanosecond, "0.000000001s")
	verifyLatency(t, testNanosAlmostSec*time.Nanosecond, "0.999999999s")
	verifyLatency(t, testLatency15MsNum*time.Millisecond, testLatency15MsStr)
	verifyLatency(t, testLatency15MsNum*time.Second, "15s")
}
