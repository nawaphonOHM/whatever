package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func createTestSpanContext(t *testing.T) context.Context {
	t.Helper()
	traceID, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("00f067aa0ba902b7")
	require.NoError(t, err)

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(context.Background(), sc)
}

func parseJSONLog(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var result map[string]any
	require.NoError(t, json.Unmarshal(data, &result))
	return result
}

func TestTraceHandler_WithActiveSpan(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))

	ctx := createTestSpanContext(t)
	logger.InfoContext(ctx, "span correlated message", "foo", "bar")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "span correlated message", logMap["msg"])
	assert.Equal(t, "bar", logMap["foo"])
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", logMap["trace_id"])
	assert.Equal(t, "00f067aa0ba902b7", logMap["span_id"])
}

func TestTraceHandler_WithStringContextKeys(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))

	ctx := context.WithValue(context.Background(), ContextKeyTraceID, "custom-trace-123")
	ctx = context.WithValue(ctx, ContextKeySpanID, "custom-span-456")

	logger.InfoContext(ctx, "custom key message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "custom key message", logMap["msg"])
	assert.Equal(t, "custom-trace-123", logMap["trace_id"])
	assert.Equal(t, "custom-span-456", logMap["span_id"])
}

func TestTraceHandler_TextHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewTextHandler(buf, nil)))

	ctx := createTestSpanContext(t)
	logger.InfoContext(ctx, "text log message", "key", "val")

	output := buf.String()
	assert.True(t, strings.Contains(output, "msg=\"text log message\""))
	assert.True(t, strings.Contains(output, "trace_id=4bf92f3577b34da6a3ce929d0e0e4736"))
	assert.True(t, strings.Contains(output, "span_id=00f067aa0ba902b7"))
}
