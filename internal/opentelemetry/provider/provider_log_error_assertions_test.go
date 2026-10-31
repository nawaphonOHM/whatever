package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/opentelemetry/config"
)

func assertDisabledProviderLogs(t *testing.T, mode providerLogMode) {
	t.Helper()
	entries, err := captureProviderLogEntries(t, &config.Config{Enabled: false}, mode)
	require.NoError(t, err)
	assertDisabledProviderEntries(t, entries, mode.levelKey)
}

func assertDisabledProviderEntries(t *testing.T, entries []map[string]any, levelKey string) {
	t.Helper()
	assertDisabledEntry(t, findEntry(entries, msgDisabled), levelKey)
	assertDisabledEntry(t, findEntry(entries, msgProvComplete), levelKey)
}

func assertDisabledEntry(t *testing.T, entry map[string]any, levelKey string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, levelInfo, entry[levelKey])
	assert.Equal(t, false, entry["enabled"])
}

func assertProviderValidationFailure(t *testing.T, mode providerLogMode) {
	t.Helper()
	cfg := &config.Config{ServiceName: "test-service", Protocol: config.ProtocolGRPC, Enabled: true}
	entries, err := captureProviderLogEntries(t, cfg, mode)
	require.Error(t, err)
	assertProviderValidationFailureEntry(t, findEntry(entries, msgValFail), mode.levelKey)
}

func assertProviderValidationFailureEntry(t *testing.T, entry map[string]any, levelKey string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, levelError, entry[levelKey])
	assert.NotEmpty(t, entry["error"])
}
