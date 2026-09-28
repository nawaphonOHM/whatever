package logging

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// ContextKey represents typed context key for trace metadata.
type ContextKey string

// Supported typed context keys for trace correlation.
const (
	ContextKeyTraceID ContextKey = "trace_id"
	ContextKeySpanID  ContextKey = "span_id"
	attrTraceID                  = "trace_id"
	attrSpanID                   = "span_id"
)

// extractSpanContext retrieves the OpenTelemetry span context from ctx.
func extractSpanContext(ctx context.Context) (trace.TraceID, trace.SpanID, bool) {
	if ctx == nil {
		return trace.TraceID{}, trace.SpanID{}, false
	}
	sc := trace.SpanFromContext(ctx).SpanContext()
	if sc.IsValid() {
		return sc.TraceID(), sc.SpanID(), true
	}
	return trace.TraceID{}, trace.SpanID{}, false
}

// matchContextValue retrieves a non-empty string value from context key.
func matchContextValue(ctx context.Context, key any) (string, bool) {
	val, ok := ctx.Value(key).(string)
	if ok && val != "" {
		return val, true
	}
	return "", false
}

// findMatchingKey iterates over keys and returns first matching string context value.
func findMatchingKey(ctx context.Context, keys []any) string {
	for _, k := range keys {
		if val, ok := matchContextValue(ctx, k); ok {
			return val
		}
	}
	return ""
}

// extractStringContextValue looks up a string value for a list of context keys.
func extractStringContextValue(ctx context.Context, keys ...any) string {
	if ctx == nil {
		return ""
	}
	return findMatchingKey(ctx, keys)
}

// appendFallbackTraceAttrs appends trace and span IDs from context string keys.
func appendFallbackTraceAttrs(ctx context.Context, record *slog.Record) {
	tid := extractStringContextValue(ctx, ContextKeyTraceID, "trace_id", "TraceID", "traceId")
	sid := extractStringContextValue(ctx, ContextKeySpanID, "span_id", "SpanID", "spanId")

	if tid != "" {
		record.AddAttrs(slog.String(attrTraceID, tid))
	}
	if sid != "" {
		record.AddAttrs(slog.String(attrSpanID, sid))
	}
}

// enrichTraceAttrs attaches trace_id and span_id attributes to the record if available.
func enrichTraceAttrs(ctx context.Context, record *slog.Record) {
	traceID, spanID, ok := extractSpanContext(ctx)
	if ok {
		record.AddAttrs(
			slog.String(attrTraceID, traceID.String()),
			slog.String(attrSpanID, spanID.String()),
		)
		return
	}
	appendFallbackTraceAttrs(ctx, record)
}
