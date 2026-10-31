package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
)

const connectionFailureTestError = "mongodb connection failed"

func TestHandleConnectionErrorLogsErrorOnlyAsAttribute(t *testing.T) {
	var output bytes.Buffer
	worker := central.DefaultWorker()
	previousLogger := worker.Logger()
	previousExit := SetExitFunc(func(int) {})
	worker.SetLogger(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() {
		worker.SetLogger(previousLogger)
		SetExitFunc(previousExit)
	})

	err := errors.New(connectionFailureTestError)
	require.ErrorIs(t, handleConnectionError(context.Background(), err), err)
	assertMongoConnectionFailureRecord(t, output.Bytes())
}

func assertMongoConnectionFailureRecord(t *testing.T, output []byte) {
	t.Helper()
	var record map[string]any
	require.NoError(t, json.Unmarshal(output, &record))
	message := "failed to connect to mongodb; exiting abnormally"
	loggedMessage, ok := record["msg"].(string)
	require.True(t, ok)
	assert.Equal(t, message, loggedMessage)
	assert.NotContains(t, loggedMessage, connectionFailureTestError)
	assert.Equal(t, connectionFailureTestError, record["error"])
}
