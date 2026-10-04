package logging_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyLogAndLogAttrsEntries(t *testing.T, entries []map[string]any) {
	t.Helper()
	require.Len(t, entries, 2)
	assert.Equal(t, testLevelWarn, entries[0]["level"])
	assert.Equal(t, "generic log warn", entries[0]["msg"])
	assert.Equal(t, "v1", entries[0]["k1"])
	assert.Equal(t, testLevelError, entries[1]["level"])
	assert.Equal(t, "generic log attrs", entries[1]["msg"])
	assert.Equal(t, "attr_v", entries[1]["attr_k"])
}

func TestLogging_LogAndLogAttrs(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureLogs(buf, logging.LevelTrace)
	defer cleanup()

	ctx := context.Background()
	logging.Log(ctx, logging.LevelWarn, "generic log warn", "k1", "v1")
	logging.LogAttrs(ctx, logging.LevelError, "generic log attrs", slog.String("attr_k", "attr_v"))

	verifyLogAndLogAttrsEntries(t, parseJSONLogs(t, buf.Bytes()))
}

func verifyFallbackEntries(t *testing.T, entries []map[string]any) {
	t.Helper()
	require.Len(t, entries, 2)
	assert.Equal(t, testLevelInfo, entries[0]["level"])
	assert.Equal(t, "fallback message", entries[0]["msg"])
	assert.Equal(t, testLevelInfo, entries[1]["level"])
	assert.Equal(t, "fallback attrs", entries[1]["msg"])
	assert.Equal(t, float64(testAnswerInt), entries[1]["n"])
}

func TestLogging_LogInvalidLevelFallback(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureLogs(buf, logging.LevelDebug)
	defer cleanup()

	ctx := context.Background()
	invalidLvl := logging.Level("INVALID_LEVEL")
	logging.Log(ctx, invalidLvl, "fallback message")
	logging.LogAttrs(ctx, invalidLvl, "fallback attrs", slog.Int("n", testAnswerInt))

	verifyFallbackEntries(t, parseJSONLogs(t, buf.Bytes()))
}

func emitNilContextLogs(ctx context.Context) {
	logging.TraceContext(ctx, "trace nil ctx")
	logging.DebugContext(ctx, "debug nil ctx")
	logging.InfoContext(ctx, "info nil ctx")
	logging.WarnContext(ctx, "warn nil ctx")
	logging.ErrorContext(ctx, "error nil ctx")
	logging.Log(ctx, logging.LevelInfo, "log nil ctx")
	logging.LogAttrs(ctx, logging.LevelInfo, "logattrs nil ctx")
}

func TestLogging_NilContextSafety(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureLogs(buf, logging.LevelTrace)
	defer cleanup()

	var nilCtx context.Context
	assert.NotPanics(t, func() {
		emitNilContextLogs(nilCtx)
	})

	assert.Len(t, parseJSONLogs(t, buf.Bytes()), testCountSeven)
}
