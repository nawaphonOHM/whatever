package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	msgLoadStart         = "loading REST configuration"
	msgLoadFallback      = "falling back to default env configuration files"
	msgLoadDone          = "REST configuration loaded successfully"
	msgLoadFail          = "failed to load REST server config"
	testExplicitPortJSON = 9099
	loaderTestFilePerm   = 0o600
)

func writeTempCustomEnv(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "custom.env")
	require.NoError(t, os.WriteFile(p, []byte(content), loaderTestFilePerm))
	return p
}

func assertEntryLevel(t *testing.T, entry map[string]any, level string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, level, entry[levelField])
}

func assertFallbackLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	assertEntryLevel(t, findLogEntry(entries, msgLoadStart), "INFO")
	fallback := findLogEntry(entries, msgLoadFallback)
	require.NotNil(t, fallback)
	assert.Equal(t, "INFO", fallback[levelField])
	assert.Equal(t, []any{".env", "configs/.env"}, fallback["files"])
	assertEntryLevel(t, findLogEntry(entries, msgLoadDone), "INFO")
}

func TestServer_LoadConfig_JSON_Fallback(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assertFallbackLogs(t, parseJSONLogs(t, buf.Bytes()))
}

func assertExplicitLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	require.NotNil(t, findLogEntry(entries, msgLoadStart))
	assert.Nil(t, findLogEntry(entries, msgLoadFallback))
	require.NotNil(t, findLogEntry(entries, msgLoadDone))
}

func TestServer_LoadConfig_JSON_ExplicitFile(t *testing.T) {
	clearServerConfigEnv(t)
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	envPath := writeTempCustomEnv(t, "OHM9996_SERVER_PORT=9099\n")
	cfg, err := LoadConfig(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, testExplicitPortJSON, cfg.Port)
	assertExplicitLogs(t, parseJSONLogs(t, buf.Bytes()))
}

func assertFailLog(t *testing.T, entries []map[string]any) {
	t.Helper()
	failEntry := findLogEntry(entries, msgLoadFail)
	require.NotNil(t, failEntry)
	assert.Equal(t, "ERROR", failEntry[levelField])
	assert.NotEmpty(t, failEntry["error"])
}

func TestServer_LoadConfig_JSON_Failure(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	t.Setenv("OHM9996_SERVER_PORT", "not-a-number")
	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assertFailLog(t, parseJSONLogs(t, buf.Bytes()))
}
