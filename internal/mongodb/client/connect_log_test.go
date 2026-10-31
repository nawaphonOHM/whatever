package client

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
)

const (
	testLogMongoPort = 27017
	msgInitClient    = "initializing mongodb client"
	msgResolveOpts   = "resolving mongodb connection options"
	msgPingSuccess   = "mongodb ping succeeded"
	msgProbeSuccess  = "mongodb probe succeeded"
	logLevelInfo     = "INFO"
)

func parseClientLogs(t *testing.T, data []byte) []map[string]any {
	t.Helper()
	var entries []map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var entry map[string]any
		require.NoError(t, decoder.Decode(&entry))
		entries = append(entries, entry)
	}
	return entries
}

func findClientEntry(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if clientEntryMatchesMessage(entry, msg) {
			return entry
		}
	}
	return nil
}

func clientEntryMatchesMessage(entry map[string]any, msg string) bool {
	return entry["msg"] == msg || entry["message"] == msg
}

func setupWorkerLogger(t *testing.T, handler slog.Handler) {
	t.Helper()
	worker := central.DefaultWorker()
	worker.SetLogger(slog.New(handler))
	t.Cleanup(func() { worker.ResetLogger() })
}
