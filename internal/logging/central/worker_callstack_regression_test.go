package central

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertCallstackRecord(t *testing.T, output string) {
	t.Helper()
	records := decodeWorkerRecords(t, output)
	record, ok := findCallstackRecord(records)
	require.True(t, ok, "callstack record not found")
	message, ok := record["msg"].(string)
	require.True(t, ok)
	assert.Equal(t, "Call stack snapshot", message)
	assert.NotContains(t, message, "testpkg.failingOperation")
	assert.Contains(t, record["callstack"], "testpkg.failingOperation")
}

func decodeWorkerRecords(t *testing.T, output string) []map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	records := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		var record map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &record))
		records = append(records, record)
	}
	return records
}

func findCallstackRecord(records []map[string]any) (map[string]any, bool) {
	for _, record := range records {
		if record["msg"] == "Call stack snapshot" {
			return record, true
		}
	}
	return nil, false
}
