package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	msgRegRoutes         = "registering REST routes"
	msgRegHealth         = "registered reserved health endpoints"
	msgRegComplete       = "route registration complete"
	regLevelInfo         = "INFO"
	expectedTotalRoutes  = 3.0
	testAppVersionConfig = "1.0.0"
)

func registerTestRoutes(t *testing.T) *Config {
	t.Helper()
	srv, err := New(&Config{AppVersion: testAppVersionConfig})
	require.NoError(t, err)
	cfg := &Config{AppVersion: testAppVersionConfig}
	require.NoError(t, srv.RegisterRoutesWithConfig(stdRegs(), cfg))
	return cfg
}

func assertRouteCounts(t *testing.T, rEntry map[string]any, lvlKey, ver string) {
	t.Helper()
	require.NotNil(t, rEntry)
	assert.Equal(t, regLevelInfo, rEntry[lvlKey])
	assert.Equal(t, float64(1), rEntry["route_count"])
	assert.Equal(t, ver, rEntry["version"])
}

func assertRegistrationAuxLogs(t *testing.T, hEntry, cEntry map[string]any, lvlKey string) {
	t.Helper()
	require.NotNil(t, hEntry)
	assert.Equal(t, regLevelInfo, hEntry[lvlKey])

	require.NotNil(t, cEntry)
	assert.Equal(t, regLevelInfo, cEntry[lvlKey])
	assert.Equal(t, expectedTotalRoutes, cEntry["total_routes"])
}

func TestServer_RegistrationLogging_JSON(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	cfg := registerTestRoutes(t)
	entries := parseJSONLogs(t, buf.Bytes())
	assertRouteCounts(t, findLogEntry(entries, msgRegRoutes), levelField, cfg.AppVersion)
	assertRegistrationAuxLogs(
		t,
		findLogEntry(entries, msgRegHealth),
		findLogEntry(entries, msgRegComplete),
		levelField,
	)
}

func TestServer_RegistrationLogging_GCP(t *testing.T) {
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	cfg := registerTestRoutes(t)
	entries := parseJSONLogs(t, buf.Bytes())
	assertRouteCounts(t, findGCPLogEntry(entries, msgRegRoutes), sevField, cfg.AppVersion)
	assertRegistrationAuxLogs(
		t,
		findGCPLogEntry(entries, msgRegHealth),
		findGCPLogEntry(entries, msgRegComplete),
		sevField,
	)
}
