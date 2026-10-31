package client

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	msgTLSFallback = "server requires TLS; attempting connection with TLS enabled"
	msgMongoFail   = "failed to connect to mongodb; exiting abnormally"
	logLevelError  = "ERROR"
)

func TestHandleConnectionError_GCP(t *testing.T) {
	assertMongoFailureLogGCP(t)
}

func assertMongoFailureLogGCP(t *testing.T) {
	t.Helper()
	entry := captureMongoFailureLogGCP(t)
	require.NotNil(t, entry)
	assert.Equal(t, logLevelError, entry["severity"])
	assert.Equal(t, "network timeout", entry["error"])
	assert.Equal(t, "ExitAbnormal", entry["exit_flag"])
}

func captureMongoFailureLogGCP(t *testing.T) map[string]any {
	t.Helper()
	var output bytes.Buffer
	setupWorkerLogger(t, core.NewGCPHandler(&output, nil))
	previousExit := SetExitFunc(func(int) {})
	t.Cleanup(func() { SetExitFunc(previousExit) })
	testErr := errors.New("network timeout")
	require.ErrorIs(t, handleConnectionError(context.Background(), testErr), testErr)
	return findClientEntry(parseClientLogs(t, output.Bytes()), msgMongoFail)
}
