package core

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

const testInvalidTraceIDNumber = 12345

func withStringKey(ctx context.Context, key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

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

	assert.Empty(t, extractStringContextValue(nilCtx, "key"))
}

func TestTraceHandler_TextHandler_FallbackKeys(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewTextHandler(buf, nil)))

	ctx := withStringKey(context.Background(), testKeyTraceID, "trace-text-999")
	ctx = withStringKey(ctx, testKeySpanID, "span-text-888")

	logger.InfoContext(ctx, "text fallback message")

	output := buf.String()
	assert.True(t, strings.Contains(output, "trace_id=trace-text-999"))
	assert.True(t, strings.Contains(output, "span_id=span-text-888"))
}

func TestFallbackTraceAttrs_AlternateKeys(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))

	ctx := withStringKey(context.Background(), "TraceID", "trace-alt-1")
	ctx = withStringKey(ctx, "SpanID", "span-alt-2")

	logger.InfoContext(ctx, "alt keys message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "trace-alt-1", logMap[testKeyTraceID])
	assert.Equal(t, "span-alt-2", logMap[testKeySpanID])
}

func TestFallbackTraceAttrs_CamelCaseAndInvalidValues(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))

	ctx := withStringKey(context.Background(), "traceId", "trace-camel-1")
	ctx = withStringKey(ctx, "spanId", "span-camel-2")
	ctx = withStringKey(ctx, testKeyTraceID, testInvalidTraceIDNumber)
	ctx = withStringKey(ctx, testKeySpanID, "")

	logger.InfoContext(ctx, "camel keys message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "trace-camel-1", logMap[testKeyTraceID])
	assert.Equal(t, "span-camel-2", logMap[testKeySpanID])
}
