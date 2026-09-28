package logger

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceIDKey     = "TraceID"
	spanIDKey      = "SpanID"
	traceIDAttrKey = "trace_id"
	spanIDAttrKey  = "span_id"
)

// extractContextString gets string value from Gin key if present.
func extractContextString(c *gin.Context, key string) string {
	val, exists := c.Get(key)
	if !exists {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

// extractKey attempts to read a string value for a key from Gin context.
func extractKey(c *gin.Context, key string) string {
	if c == nil {
		return ""
	}
	return extractContextString(c, key)
}

// spanContext retrieves the span context from the Gin request context.
func spanContext(c *gin.Context) trace.SpanContext {
	if c == nil || c.Request == nil {
		return trace.SpanContext{}
	}
	return trace.SpanFromContext(c.Request.Context()).SpanContext()
}

// extractSpanTraceID retrieves trace ID from active request span context.
func extractSpanTraceID(c *gin.Context) string {
	sc := spanContext(c)
	if !sc.IsValid() {
		return ""
	}
	return sc.TraceID().String()
}

// extractSpanContextID retrieves span ID from active request span context.
func extractSpanContextID(c *gin.Context) string {
	sc := spanContext(c)
	if !sc.IsValid() {
		return ""
	}
	return sc.SpanID().String()
}

// extractTraceID gets the trace ID from the Gin context keys or active request span.
func extractTraceID(c *gin.Context) string {
	if id := extractKey(c, traceIDKey); id != "" {
		return id
	}
	if id := extractKey(c, traceIDAttrKey); id != "" {
		return id
	}
	return extractSpanTraceID(c)
}

// extractSpanID gets the span ID from the Gin context keys or active request span.
func extractSpanID(c *gin.Context) string {
	if id := extractKey(c, spanIDKey); id != "" {
		return id
	}
	if id := extractKey(c, spanIDAttrKey); id != "" {
		return id
	}
	return extractSpanContextID(c)
}

// traceAttrs returns log attributes for trace_id and span_id if available.
func traceAttrs(c *gin.Context) []slog.Attr {
	var list []slog.Attr
	if traceID := extractTraceID(c); traceID != "" {
		list = append(list, slog.String(traceIDAttrKey, traceID))
	}
	if spanID := extractSpanID(c); spanID != "" {
		list = append(list, slog.String(spanIDAttrKey, spanID))
	}
	return list
}
