package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestExtractSpanContext_NilOrInvalid(t *testing.T) {
	var nilCtx context.Context
	tid, sid, ok := extractSpanContext(nilCtx)
	assert.False(t, ok)
	assert.Equal(t, trace.TraceID{}, tid)
	assert.Equal(t, trace.SpanID{}, sid)

	tid, sid, ok = extractSpanContext(context.Background())
	assert.False(t, ok)
	assert.Equal(t, trace.TraceID{}, tid)
	assert.Equal(t, trace.SpanID{}, sid)
}

//nolint:revive,staticcheck // Testing raw string context keys for fallback compatibility.
func TestTraceHandler_TextHandler_FallbackKeys(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewTextHandler(buf, nil)))

	ctx := context.WithValue(context.Background(), "trace_id", "trace-text-999")
	ctx = context.WithValue(ctx, "span_id", "span-text-888")

	logger.InfoContext(ctx, "text fallback message")

	output := buf.String()
	assert.True(t, strings.Contains(output, "trace_id=trace-text-999"))
	assert.True(t, strings.Contains(output, "span_id=span-text-888"))
}

//nolint:revive,staticcheck // Testing raw string context keys for fallback compatibility.
func TestFallbackTraceAttrs_AlternateKeys(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))

	ctx := context.WithValue(context.Background(), "TraceID", "trace-alt-1")
	ctx = context.WithValue(ctx, "SpanID", "span-alt-2")

	logger.InfoContext(ctx, "alt keys message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "trace-alt-1", logMap["trace_id"])
	assert.Equal(t, "span-alt-2", logMap["span_id"])
}

//nolint:revive,staticcheck // Testing raw string context keys for fallback compatibility.
func TestFallbackTraceAttrs_CamelCaseAndInvalidValues(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))

	ctx := context.WithValue(context.Background(), "traceId", "trace-camel-1")
	ctx = context.WithValue(ctx, "spanId", "span-camel-2")
	ctx = context.WithValue(ctx, "trace_id", 12345)
	ctx = context.WithValue(ctx, "span_id", "")

	logger.InfoContext(ctx, "camel keys message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "trace-camel-1", logMap["trace_id"])
	assert.Equal(t, "span-camel-2", logMap["span_id"])
}
