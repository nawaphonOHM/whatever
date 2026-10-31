package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/opentelemetry/config"
)

func assertActiveProviderLogs(t *testing.T, mode providerLogMode) {
	t.Helper()
	entries, err := captureProviderLogEntries(t, activeProviderTestConfig(), mode)
	require.NoError(t, err)
	assertActiveProviderEntries(t, entries, mode.levelKey)
}

func assertActiveProviderEntries(t *testing.T, entries []map[string]any, levelKey string) {
	t.Helper()
	choices := findEntry(entries, msgProvChoices)
	exporter := findEntry(entries, msgExpChoice)
	complete := findEntry(entries, msgProvComplete)
	require.NotNil(t, choices)
	require.NotNil(t, exporter)
	require.NotNil(t, complete)
	assertProviderChoiceAttrs(t, choices, levelKey)
	assertExporterChoiceAttrs(t, exporter, levelKey)
	assertProviderCompletionAttrs(t, complete, levelKey)
}

func assertProviderChoiceAttrs(t *testing.T, entry map[string]any, levelKey string) {
	t.Helper()
	assert.Equal(t, levelInfo, entry[levelKey])
	assert.Equal(t, "test-service", entry["service_name"])
	assert.Equal(t, "127.0.0.1:4317", entry["endpoint"])
	assert.Equal(t, config.ProtocolGRPC, entry["protocol"])
	assert.Equal(t, testSampleRateHalf, entry["sample_rate"])
}

func assertExporterChoiceAttrs(t *testing.T, entry map[string]any, levelKey string) {
	t.Helper()
	assert.Equal(t, levelInfo, entry[levelKey])
	assert.Equal(t, config.ProtocolGRPC, entry["protocol"])
	assert.Equal(t, "grpc", entry["transport"])
	assert.Equal(t, "127.0.0.1:4317", entry["endpoint"])
}

func assertProviderCompletionAttrs(t *testing.T, entry map[string]any, levelKey string) {
	t.Helper()
	assert.Equal(t, levelInfo, entry[levelKey])
	assert.Equal(t, true, entry["enabled"])
}
