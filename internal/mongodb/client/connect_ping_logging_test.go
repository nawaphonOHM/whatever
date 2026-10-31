package client

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

const (
	clientLogFormatJSON = "json"
	clientLogFormatGCP  = "gcp"
)

type clientLogMode struct {
	name     string
	levelKey string
	format   string
}

func clientLogModes() []clientLogMode {
	return []clientLogMode{
		{name: clientLogFormatJSON, levelKey: "level", format: clientLogFormatJSON},
		{name: clientLogFormatGCP, levelKey: "severity", format: clientLogFormatGCP},
	}
}

func (m clientLogMode) handler(output io.Writer) slog.Handler {
	if m.format == clientLogFormatGCP {
		return core.NewGCPHandler(output, nil)
	}
	return slog.NewJSONHandler(output, nil)
}

func TestVerifyPingAndOptsLog_JSONAndGCP(t *testing.T) {
	for _, mode := range clientLogModes() {
		t.Run(mode.name, func(t *testing.T) { assertPingAndOptionsLog(t, mode) })
	}
}

func assertPingAndOptionsLog(t *testing.T, mode clientLogMode) {
	t.Helper()
	var output bytes.Buffer
	setupWorkerLogger(t, mode.handler(&output))
	restore := SetMockPing(func(context.Context, *mongo.Client) error { return nil })
	t.Cleanup(restore)
	cfg := pingLoggingTestConfig()
	logConnectionOptions(context.Background(), cfg, true)
	require.NoError(t, verifyPing(context.Background(), cfg, nil))
	assertPingAndOptionsEntries(t, parseClientLogs(t, output.Bytes()), mode)
}

func pingLoggingTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.AuthSource = "admin"
	return cfg
}

func assertPingAndOptionsEntries(t *testing.T, entries []map[string]any, mode clientLogMode) {
	t.Helper()
	assertOptionsEntry(t, findClientEntry(entries, msgResolveOpts), mode)
	assertPingEntry(t, findClientEntry(entries, msgPingSuccess), mode)
}

func assertOptionsEntry(t *testing.T, entry map[string]any, mode clientLogMode) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, logLevelInfo, entry[mode.levelKey])
	assert.Equal(t, true, entry["tls"])
	assert.Equal(t, "admin", entry["auth_source"])
}

func assertPingEntry(t *testing.T, entry map[string]any, mode clientLogMode) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, logLevelInfo, entry[mode.levelKey])
}
