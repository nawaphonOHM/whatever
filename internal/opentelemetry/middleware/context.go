package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
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

// extractSpanTraceID retrieves trace ID from active request span context.
func extractSpanTraceID(c *gin.Context) string {
	if isInvalidContext(c) {
		return ""
	}
	sc := trace.SpanFromContext(c.Request.Context()).SpanContext()
	if !sc.IsValid() {
		return ""
	}
	return sc.TraceID().String()
}

// extractSpanID retrieves span ID from active request span context.
func extractSpanID(c *gin.Context) string {
	if isInvalidContext(c) {
		return ""
	}
	sc := trace.SpanFromContext(c.Request.Context()).SpanContext()
	if !sc.IsValid() {
		return ""
	}
	return sc.SpanID().String()
}

// GetTraceID extracts the trace ID from the Gin context or active request span.
func GetTraceID(c *gin.Context) string {
	if id := extractKey(c, TraceIDKey); id != "" {
		return id
	}
	if id := extractKey(c, TraceIDAttrKey); id != "" {
		return id
	}
	return extractSpanTraceID(c)
}

// GetSpanID extracts the span ID from the Gin context or active request span.
func GetSpanID(c *gin.Context) string {
	if id := extractKey(c, SpanIDKey); id != "" {
		return id
	}
	if id := extractKey(c, SpanIDAttrKey); id != "" {
		return id
	}
	return extractSpanID(c)
}

// applySpanContextKeys stores valid span context identifiers into context.
func applySpanContextKeys(c *gin.Context, sc trace.SpanContext) {
	if !sc.IsValid() {
		return
	}
	c.Set(TraceIDKey, sc.TraceID().String())
	c.Set(TraceIDAttrKey, sc.TraceID().String())
	c.Set(SpanIDKey, sc.SpanID().String())
	c.Set(SpanIDAttrKey, sc.SpanID().String())
}

// injectTraceContext stores trace and span IDs in Gin context keys.
func injectTraceContext(c *gin.Context, span trace.Span) {
	if c == nil || span == nil {
		return
	}
	applySpanContextKeys(c, span.SpanContext())
}
