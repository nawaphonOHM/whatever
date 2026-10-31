package client

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

func TestValidateAndLogConfig_JSONAndGCP(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = testLogMongoPort
	cfg.Database = "testdb"
	cfg.AuthSource = "admin"
	cfg.UUIDRepresentation = config.UUIDRepresentationStandard
	assertConfigLogJSON(t, cfg)
	assertConfigLogGCP(t, cfg)
}

func assertConfigLogJSON(t *testing.T, cfg *config.Config) {
	t.Helper()
	var output bytes.Buffer
	setupWorkerLogger(t, slog.NewJSONHandler(&output, nil))
	require.NoError(t, validateAndLogConfig(context.Background(), cfg))
	entry := findClientEntry(parseClientLogs(t, output.Bytes()), msgInitClient)
	require.NotNil(t, entry)
	assert.Equal(t, logLevelInfo, entry["level"])
	assertConfigLogAttributes(t, entry)
}

func assertConfigLogGCP(t *testing.T, cfg *config.Config) {
	t.Helper()
	var output bytes.Buffer
	setupWorkerLogger(t, core.NewGCPHandler(&output, nil))
	require.NoError(t, validateAndLogConfig(context.Background(), cfg))
	entry := findClientEntry(parseClientLogs(t, output.Bytes()), msgInitClient)
	require.NotNil(t, entry)
	assert.Equal(t, logLevelInfo, entry["severity"])
	assertConfigLogAttributes(t, entry)
}

func assertConfigLogAttributes(t *testing.T, entry map[string]any) {
	t.Helper()
	assert.Equal(t, "127.0.0.1", entry["host"])
	assert.Equal(t, float64(testLogMongoPort), entry["port"])
	assert.Equal(t, "testdb", entry["database"])
}
