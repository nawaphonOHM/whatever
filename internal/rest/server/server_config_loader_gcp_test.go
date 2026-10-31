package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testExplicitPortGCP = 9098

func assertGCPEntryLevel(t *testing.T, entry map[string]any, level string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, level, entry[sevField])
}

func assertGCPFallbackLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	assertGCPEntryLevel(t, findGCPLogEntry(entries, msgLoadStart), gcpInfoLevel)
	fallback := findGCPLogEntry(entries, msgLoadFallback)
	require.NotNil(t, fallback)
	assert.Equal(t, gcpInfoLevel, fallback[sevField])
	assert.Equal(t, []any{".env", "configs/.env"}, fallback["files"])
	assertGCPEntryLevel(t, findGCPLogEntry(entries, msgLoadDone), gcpInfoLevel)
}

func TestServer_LoadConfig_GCP_Fallback(t *testing.T) {
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assertGCPFallbackLogs(t, parseJSONLogs(t, buf.Bytes()))
}

func assertGCPExplicitLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	require.NotNil(t, findGCPLogEntry(entries, msgLoadStart))
	assert.Nil(t, findGCPLogEntry(entries, msgLoadFallback))
	require.NotNil(t, findGCPLogEntry(entries, msgLoadDone))
}

func TestServer_LoadConfig_GCP_ExplicitFile(t *testing.T) {
	clearServerConfigEnv(t)
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	envPath := writeTempCustomEnv(t, "OHM9996_SERVER_PORT=9098\n")
	cfg, err := LoadConfig(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, testExplicitPortGCP, cfg.Port)
	assertGCPExplicitLogs(t, parseJSONLogs(t, buf.Bytes()))
}

func assertGCPFailLog(t *testing.T, entries []map[string]any) {
	t.Helper()
	failEntry := findGCPLogEntry(entries, msgLoadFail)
	require.NotNil(t, failEntry)
	assert.Equal(t, "ERROR", failEntry[sevField])
	assert.NotEmpty(t, failEntry["error"])
}

func TestServer_LoadConfig_GCP_Failure(t *testing.T) {
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	t.Setenv("OHM9996_SERVER_PORT", "invalid-port")
	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assertGCPFailLog(t, parseJSONLogs(t, buf.Bytes()))
}
