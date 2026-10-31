package client

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

const (
	msgTLSFallback = "server requires TLS; attempting connection with TLS enabled"
	msgMongoFail   = "failed to connect to mongodb; exiting abnormally"
	msgInitSuccess = "mongodb client initialization complete"
	logLevelError  = "ERROR"
)

func TestHandleConnectionError_GCP(t *testing.T) {
	assertMongoFailureLogGCP(t)
}

func assertMongoFailureLogGCP(t *testing.T) {
	t.Helper()
	entry := captureMongoFailureLogGCP(t)
	require.NotNil(t, entry)
	assertGCPTimestamp(t, entry)
	assertGCPSeverity(t, entry, logLevelError)
	assert.Equal(t, msgMongoFail, entry["message"])
	assert.Equal(t, "network timeout", entry["error"])
	assert.Equal(t, "ExitAbnormal", entry["exit_flag"])
}

func captureMongoFailureLogGCP(t *testing.T) map[string]any {
	t.Helper()
	var output bytes.Buffer
	setupGCPTestLogger(t, &output)
	previousExit := SetExitFunc(func(int) {})
	t.Cleanup(func() { SetExitFunc(previousExit) })
	testErr := errors.New("network timeout")
	require.ErrorIs(t, handleConnectionError(context.Background(), testErr), testErr)
	return findClientEntry(parseClientLogs(t, output.Bytes()), msgMongoFail)
}

func captureMongoLifecycleLogsGCP(t *testing.T, cfg *config.Config) []map[string]any {
	t.Helper()
	var buf bytes.Buffer
	setupGCPTestLogger(t, &buf)
	stubMongoPingAndProbe(t)
	_, err := ConnectWithConfig(context.Background(), cfg)
	require.NoError(t, err)
	return parseClientLogs(t, buf.Bytes())
}

func TestMongoDBConnectLifecycleLogs_GCP(t *testing.T) {
	entries := captureMongoLifecycleLogsGCP(t, gcpMongoTestConfig())
	assertMongoInitEntryGCP(t, findClientEntry(entries, msgInitClient))
	assertMongoOptsEntryGCP(t, findClientEntry(entries, msgResolveOpts))
	assertMongoPingEntryGCP(t, findClientEntry(entries, msgPingSuccess))
	assertMongoCompleteEntryGCP(t, findClientEntry(entries, msgInitSuccess))
}
