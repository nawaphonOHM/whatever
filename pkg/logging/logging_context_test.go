package logging_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func createContextWithSpan(t *testing.T) context.Context {
	t.Helper()
	traceID, err := trace.TraceIDFromHex(testTraceIDHex)
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex(testSpanIDHex)
	require.NoError(t, err)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(context.Background(), spanCtx)
}

func emitContextLogs(ctx context.Context) {
	logging.TraceContext(ctx, "trace ctx", "t_ctx", testNumOne)
	logging.DebugContext(ctx, "debug ctx", "d_ctx", testNumTwo)
	logging.InfoContext(ctx, "info ctx", "i_ctx", testNumThree)
	logging.WarnContext(ctx, "warn ctx", "w_ctx", testNumFour)
	logging.ErrorContext(ctx, "error ctx", "e_ctx", testNumFive)
}

func verifyContextEntries(t *testing.T, entries []map[string]any) {
	require.Len(t, entries, testCountFive)
	for _, entry := range entries {
		assert.Equal(t, testTraceIDHex, entry[testKeyTraceID])
		assert.Equal(t, testSpanIDHex, entry[testKeySpanID])
	}
	assert.Equal(t, testLevelTrace, entries[0][testKeyLevel])
	assert.Equal(t, testLevelDebug, entries[1][testKeyLevel])
	assert.Equal(t, testLevelInfo, entries[2][testKeyLevel])
	assert.Equal(t, testLevelWarn, entries[testIndexThree][testKeyLevel])
	assert.Equal(t, testLevelError, entries[testIndexFour][testKeyLevel])
}

func TestLogging_ContextFunctions(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureLogs(buf, logging.LevelTrace)
	defer cleanup()

	ctx := createContextWithSpan(t)
	emitContextLogs(ctx)
	verifyContextEntries(t, parseJSONLogs(t, buf.Bytes()))
}
