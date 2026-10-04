package logging_test

import (
	"bytes"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevel_Constants(t *testing.T) {
	tests := []struct {
		name     string
		level    logging.Level
		expected string
	}{
		{"Trace", logging.LevelTrace, testLevelTrace},
		{"Debug", logging.LevelDebug, testLevelDebug},
		{"Info", logging.LevelInfo, testLevelInfo},
		{"Warn", logging.LevelWarn, testLevelWarn},
		{"Error", logging.LevelError, testLevelError},
		{"Fatal", logging.LevelFatal, testLevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.level))
			assert.Equal(t, tt.expected, tt.level.String())
			slogLvl, err := tt.level.SlogLevel()
			assert.NoError(t, err)
			assert.NotEmpty(t, slogLvl.String())
		})
	}
}

func emitSeverityLogs() {
	logging.Trace("trace message", "key_t", "val_t")
	logging.Debug("debug message", "key_d", "val_d")
	logging.Info("info message", "key_i", "val_i")
	logging.Warn("warn message", "key_w", "val_w")
	logging.Error("error message", "key_e", "val_e")
}

func verifyTraceDebugInfo(t *testing.T, entries []map[string]any) {
	t.Helper()
	assert.Equal(t, testLevelTrace, entries[0]["level"])
	assert.Equal(t, "trace message", entries[0]["msg"])
	assert.Equal(t, "val_t", entries[0]["key_t"])

	assert.Equal(t, testLevelDebug, entries[1]["level"])
	assert.Equal(t, "debug message", entries[1]["msg"])
	assert.Equal(t, "val_d", entries[1]["key_d"])

	assert.Equal(t, testLevelInfo, entries[2]["level"])
	assert.Equal(t, "info message", entries[2]["msg"])
	assert.Equal(t, "val_i", entries[2]["key_i"])
}

func verifyWarnError(t *testing.T, entries []map[string]any) {
	t.Helper()
	assert.Equal(t, testLevelWarn, entries[3]["level"])
	assert.Equal(t, "warn message", entries[3]["msg"])
	assert.Equal(t, "val_w", entries[3]["key_w"])

	assert.Equal(t, testLevelError, entries[4]["level"])
	assert.Equal(t, "error message", entries[4]["msg"])
	assert.Equal(t, "val_e", entries[4]["key_e"])
}

func verifySeverityEntries(t *testing.T, entries []map[string]any) {
	require.Len(t, entries, testCountFive)
	verifyTraceDebugInfo(t, entries)
	verifyWarnError(t, entries)
}

func TestLogging_SeverityFunctions(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureLogs(buf, logging.LevelTrace)
	defer cleanup()

	emitSeverityLogs()
	verifySeverityEntries(t, parseJSONLogs(t, buf.Bytes()))
}

func TestLogging_Flush(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureLogs(buf, logging.LevelInfo)
	defer cleanup()

	logging.Info("test flush message", "key_f", "val_f")
	logging.Flush()

	entries := parseJSONLogs(t, buf.Bytes())
	require.Len(t, entries, 1)
	assert.Equal(t, "test flush message", entries[0]["msg"])
	assert.Equal(t, "val_f", entries[0]["key_f"])
}
