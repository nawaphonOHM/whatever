package core_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func normalizeCoreKeys(m map[string]any) {
	if ts, ok := m["timestamp"]; ok {
		m[slog.TimeKey] = ts
	}
	if sev, ok := m["severity"]; ok {
		m[slog.LevelKey] = sev
	}
}

func normalizeSlogtestKeys(m map[string]any) {
	normalizeCoreKeys(m)
	if msg, ok := m["message"]; ok {
		m[slog.MessageKey] = msg
	}
}

func parseSingleLine(line string) (map[string]any, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil, false
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(trimmed), &m); err != nil {
		return nil, false
	}
	normalizeSlogtestKeys(m)
	return m, true
}

func parseResultsBuffer(buf *bytes.Buffer) []map[string]any {
	var entries []map[string]any
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	buf.Reset()
	for _, line := range lines {
		if m, ok := parseSingleLine(line); ok {
			entries = append(entries, m)
		}
	}
	return entries
}

func TestGCPHandler_Slogtest(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	results := func() []map[string]any {
		return parseResultsBuffer(&buf)
	}

	err := slogtest.TestHandler(handler, results)
	require.NoError(t, err)
}
