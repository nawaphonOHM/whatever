package logging_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTraceIDHex      = "4bf92f3577b34da6a3ce929d0e0e4736"
	testSpanIDHex       = "00f067aa0ba902b7"
	testFatalTraceIDHex = "1234567890abcdef1234567890abcdef"
	testFatalSpanIDHex  = "abcdef1234567890"

	testPort      = 8080
	testLatencyMs = 250
	testAnswerInt = 42

	testCountFive  = 5
	testCountSeven = 7
	testNumOne     = 1
	testNumTwo     = 2
	testNumThree   = 3
	testNumFour    = 4
	testNumFive    = 5

	testIndexThree = 3
	testIndexFour  = 4

	testKeyLevel   = "level"
	testKeyMsg     = "msg"
	testKeyTraceID = "trace_id"
	testKeySpanID  = "span_id"

	testLevelTrace = "TRACE"
	testLevelDebug = "DEBUG"
	testLevelInfo  = "INFO"
	testLevelWarn  = "WARN"
	testLevelError = "ERROR"
	testLevelFatal = "FATAL"
)

func parseJSONLog(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal(data, &entry))
	return entry
}

func parseJSONLogs(t *testing.T, data []byte) []map[string]any {
	t.Helper()
	var entries []map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var entry map[string]any
		require.NoError(t, decoder.Decode(&entry))
		entries = append(entries, entry)
	}
	return entries
}

func verifyTextOutput(t *testing.T, output string) {
	t.Helper()
	assert.Contains(t, output, "level=INFO")
	assert.Contains(t, output, "msg=\"server starting\"")
	assert.Contains(t, output, "port=8080")
	assert.Contains(t, output, "level=WARN")
	assert.Contains(t, output, "msg=\"high latency detected\"")
	assert.Contains(t, output, "ms=250")
}

func TestLogging_TextOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureTextLogs(buf, logging.LevelInfo)
	defer cleanup()

	logging.Info("server starting", "port", testPort)
	logging.Warn("high latency detected", "ms", testLatencyMs)

	verifyTextOutput(t, buf.String())
}
