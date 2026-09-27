package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// extractURI returns the non-empty URI or fallback path.
func extractURI(c *gin.Context) string {
	if uri := c.Request.URL.RequestURI(); uri != "" {
		return uri
	}
	return c.Request.URL.Path
}

// getTargetURI returns the full request URI or fallback path safely.
func getTargetURI(c *gin.Context) string {
	if isInvalidContext(c) {
		return ""
	}
	if c.Request.URL == nil {
		return ""
	}
	return extractURI(c)
}

// buildBaseAttributes creates standard method, target, and client IP attributes.
func buildBaseAttributes(c *gin.Context) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String(AttrHTTPMethod, c.Request.Method),
		attribute.String(AttrHTTPTarget, getTargetURI(c)),
		attribute.String(AttrClientAddress, c.ClientIP()),
	}
}

// buildInitialAttributes collects request-level semantic attributes.
func buildInitialAttributes(c *gin.Context) []attribute.KeyValue {
	if isInvalidContext(c) {
		return nil
	}
	attrs := buildBaseAttributes(c)
	if route := c.FullPath(); route != "" {
		attrs = append(attrs, attribute.String(AttrHTTPRoute, route))
	}
	return attrs
}

// buildInitialSpanName generates the provisional span name before routing.
func buildInitialSpanName(c *gin.Context) string {
	if isInvalidContext(c) {
		return "HTTP"
	}
	if route := c.FullPath(); route != "" {
		return c.Request.Method + " " + route
	}
	return c.Request.Method
}

// updateSpanRoute updates the span name and route attribute when full path is matched.
func updateSpanRoute(span trace.Span, c *gin.Context) {
	if route := c.FullPath(); route != "" {
		span.SetName(c.Request.Method + " " + route)
		span.SetAttributes(attribute.String(AttrHTTPRoute, route))
	}
}

// enrichSpanAttributes records status and route attributes after next.
func enrichSpanAttributes(span trace.Span, c *gin.Context) {
	updateSpanRoute(span, c)
	statusCode := c.Writer.Status()
	span.SetAttributes(attribute.Int(AttrHTTPStatusCode, statusCode))
	recordSpanStatus(span, c, statusCode)
}

// enrichSpanAfterNext updates span name, route, and status code after execution.
func enrichSpanAfterNext(span trace.Span, c *gin.Context) {
	if span == nil {
		return
	}
	if isInvalidContext(c) {
		return
	}
	enrichSpanAttributes(span, c)
}
