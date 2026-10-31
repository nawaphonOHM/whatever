package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	msgChoices  = "REST server configuration choices"
	msgComplete = "REST server initialization complete"
	choicesInfo = "INFO"
	levelField  = "level"
	sevField    = "severity"
)

func assertChoicesLog(t *testing.T, entry map[string]any, lvlKey string, cfg *Config) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, choicesInfo, entry[lvlKey])
	assert.Equal(t, cfg.Host, entry["host"])
	assert.Equal(t, float64(cfg.Port), entry["port"])
}

func assertCompleteLog(t *testing.T, entry map[string]any, lvlKey string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, choicesInfo, entry[lvlKey])
	assert.NotEmpty(t, entry["addr"])
}

func TestServer_ConfigChoicesLogging_JSON(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	cfg := DefaultConfig()
	srv, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	entries := parseJSONLogs(t, buf.Bytes())
	assertChoicesLog(t, findLogEntry(entries, msgChoices), levelField, cfg)
	assertCompleteLog(t, findLogEntry(entries, msgComplete), levelField)
}

func TestServer_ConfigChoicesLogging_GCP(t *testing.T) {
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	cfg := DefaultConfig()
	srv, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	entries := parseJSONLogs(t, buf.Bytes())
	assertChoicesLog(t, findGCPLogEntry(entries, msgChoices), sevField, cfg)
	assertCompleteLog(t, findGCPLogEntry(entries, msgComplete), sevField)
}
