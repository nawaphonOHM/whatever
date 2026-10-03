package logging_test

import (
	"bytes"
	"encoding/json"
	"testing"

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
