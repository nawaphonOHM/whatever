package provider

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	msgProvChoices  = "OpenTelemetry provider configuration choices"
	msgExpChoice    = "OpenTelemetry exporter choice"
	msgProvComplete = "OpenTelemetry provider initialization complete"
	msgDisabled     = "OpenTelemetry disabled; falling back to noop tracer provider"
	fieldLevel      = "level"
	levelInfo       = "INFO"
	fieldSeverity   = "severity"
	levelError      = "ERROR"
	msgValFail      = "OpenTelemetry provider configuration validation failed"
)

func parseLogEntries(t *testing.T, data []byte) []map[string]any {
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

func findEntry(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if providerEntryMatchesMessage(entry, msg) {
			return entry
		}
	}
	return nil
}

func providerEntryMatchesMessage(entry map[string]any, msg string) bool {
	return entry["msg"] == msg || entry["message"] == msg
}

func TestInitTracerProvider_JSONLogging_Active(t *testing.T) {
	assertActiveProviderLogs(t, jsonProviderLogMode())
}
