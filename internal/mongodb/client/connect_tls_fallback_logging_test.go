package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

func TestTLSFallbackLogsConnectionTarget_JSONAndGCP(t *testing.T) {
	for _, mode := range clientLogModes() {
		t.Run(mode.name, func(t *testing.T) { assertTLSFallbackLog(t, mode) })
	}
}

func assertTLSFallbackLog(t *testing.T, mode clientLogMode) {
	t.Helper()
	var output bytes.Buffer
	setupWorkerLogger(t, mode.handler(&output))
	runFallbackTLSAttempt(t)
	entry := findClientEntry(parseClientLogs(t, output.Bytes()), msgTLSFallback)
	require.NotNil(t, entry)
	assertTLSFallbackAttributes(t, entry, mode)
}

func runFallbackTLSAttempt(t *testing.T) {
	t.Helper()
	setFallbackPingAndProbe(t)
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = testLogMongoPort
	cfg.AuthSource = "admin"
	cfg.UUIDRepresentation = config.UUIDRepresentationStandard
	connected, err := fallbackTLSAttempt(context.Background(), cfg)
	require.NoError(t, err)
	require.NoError(t, connected.Disconnect(context.Background()))
}

func assertTLSFallbackAttributes(t *testing.T, entry map[string]any, mode clientLogMode) {
	t.Helper()
	assert.Equal(t, logLevelInfo, entry[mode.levelKey])
	assert.Equal(t, "127.0.0.1", entry["host"])
	assert.Equal(t, float64(testLogMongoPort), entry["port"])
}

func setFallbackPingAndProbe(t *testing.T) {
	t.Helper()
	cleanupPing := SetMockPing(func(context.Context, *mongo.Client) error { return nil })
	cleanupProbe := SetMockProbe(func(context.Context, *Client, string) error { return nil })
	t.Cleanup(cleanupPing)
	t.Cleanup(cleanupProbe)
}
