package config

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	logMsgChoices  = "OpenTelemetry configuration choices"
	logMsgComplete = "OpenTelemetry configuration initialization complete"
	fieldSeverity  = "severity"
	fieldLevel     = "level"
	levelInfo      = "INFO"
)

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

func findLogEntry(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if logEntryMatchesMessage(entry, msg) {
			return entry
		}
	}
	return nil
}

func logEntryMatchesMessage(entry map[string]any, msg string) bool {
	return entry["msg"] == msg || entry["message"] == msg
}

func assertChoicesAttrs(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, DefaultEnabled, entry["enabled"])
	assert.Equal(t, DefaultServiceName, entry["service_name"])
	assert.Equal(t, DefaultEndpoint, entry["endpoint"])
	assert.Equal(t, DefaultProtocol, entry["protocol"])
	assert.Equal(t, DefaultInsecure, entry["insecure"])
	assert.Equal(t, DefaultSampleRate, entry["sample_rate"])
}

func TestLoad_JSONLogging(t *testing.T) {
	assertConfigLoadFormat(t, fieldLevel, captureJSONLoadLogs)
}

func TestLoad_GCPLogging(t *testing.T) {
	assertConfigLoadFormat(t, fieldSeverity, captureGCPLoadLogs)
}

func captureJSONLoadLogs(output *bytes.Buffer) func() {
	return core.CaptureLogs(output)
}

func captureGCPLoadLogs(output *bytes.Buffer) func() {
	handler := core.NewGCPHandler(output, nil)
	return core.SetTestLogger(core.NewWithHandler(handler))
}

func assertConfigLoadFormat(t *testing.T, levelKey string, capture func(*bytes.Buffer) func()) {
	entries := captureConfigLoadLogs(t, capture)
	choices := findLogEntry(entries, logMsgChoices)
	assertChoicesAttrs(t, choices)
	assertLogLevel(t, choices, levelKey)
	complete := findLogEntry(entries, logMsgComplete)
	assertLogLevel(t, complete, levelKey)
}

func captureConfigLoadLogs(t *testing.T, capture func(*bytes.Buffer) func()) []map[string]any {
	t.Helper()
	clearEnv(t)
	output := new(bytes.Buffer)
	reset := capture(output)
	defer reset()
	_, err := Load()
	require.NoError(t, err)
	return parseJSONLogs(t, output.Bytes())
}

func assertLogLevel(t *testing.T, entry map[string]any, field string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, levelInfo, entry[field])
}
