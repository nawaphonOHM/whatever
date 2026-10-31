package client

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
)

func setStorageTestLogger(t *testing.T, handler slog.Handler) {
	t.Helper()
	worker := central.DefaultWorker()
	worker.SetLogger(slog.New(handler))
	t.Cleanup(func() { worker.ResetLogger() })
}

func parseStorageLogs(t *testing.T, data []byte) []map[string]any {
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

func findStorageLog(entries []map[string]any, message string) map[string]any {
	for _, entry := range entries {
		if storageLogMatchesMessage(entry, message) {
			return entry
		}
	}
	return nil
}

func storageLogMatchesMessage(entry map[string]any, message string) bool {
	return entry["msg"] == message || entry["message"] == message
}
