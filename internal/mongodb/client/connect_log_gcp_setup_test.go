package client

import (
	"bytes"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

var validGCPSeverities = map[string]struct{}{
	"DEBUG":     {},
	"INFO":      {},
	"NOTICE":    {},
	"WARNING":   {},
	"ERROR":     {},
	"CRITICAL":  {},
	"ALERT":     {},
	"EMERGENCY": {},
}

func assertDefaultWorkerFollowsSlogDefault(t *testing.T) {
	t.Helper()
	var canaryBuf bytes.Buffer
	canaryLogger := slog.New(slog.NewTextHandler(&canaryBuf, nil))
	orig := slog.Default()
	slog.SetDefault(canaryLogger)
	msg := "default worker must dynamically track slog.Default()"
	assert.Equal(t, canaryLogger, central.DefaultWorker().Logger(), msg)
	slog.SetDefault(orig)
}

func setupGCPTestLogger(t *testing.T, buf *bytes.Buffer) {
	t.Helper()
	origDefault, origWriter, origFlags := slog.Default(), log.Writer(), log.Flags()
	central.DefaultWorker().ResetLogger()
	slog.SetDefault(slog.New(core.NewGCPHandler(buf, nil)))
	log.SetOutput(buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		slog.SetDefault(origDefault)
		log.SetOutput(origWriter)
		log.SetFlags(origFlags)
		central.DefaultWorker().ResetLogger()
		assertDefaultWorkerFollowsSlogDefault(t)
	})
}

func assertGCPSeverity(t *testing.T, entry map[string]any, expected string) {
	t.Helper()
	sev, ok := entry[logFieldSeverity].(string)
	require.True(t, ok, "expected severity string, got %T: %v", entry[logFieldSeverity], entry[logFieldSeverity])
	assert.NotEqual(t, "DEFAULT", sev, "GCP severity must not be DEFAULT")
	assert.Equal(t, expected, sev)
	_, valid := validGCPSeverities[sev]
	assert.True(t, valid, "severity %q is not in standard GCP severity enum", sev)
}

func assertGCPTimestamp(t *testing.T, entry map[string]any) {
	t.Helper()
	ts, ok := entry["timestamp"].(string)
	require.True(t, ok, "expected timestamp string, got %T: %v", entry["timestamp"], entry["timestamp"])
	parsed, err := time.Parse(time.RFC3339Nano, ts)
	require.NoError(t, err, "timestamp %q is not valid RFC3339", ts)
	assert.Equal(t, time.UTC, parsed.Location(), "timestamp location must be UTC")
}
