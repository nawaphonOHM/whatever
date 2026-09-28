// Package middleware provides Gin HTTP request tracing and W3C context propagation.
package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// isTraceDisabled returns true if configuration is nil or telemetry is disabled.
func isTraceDisabled(cfg *config.Config) bool {
	return cfg == nil || !cfg.Enabled
}

// isInvalidContext returns true if Gin context or HTTP request is nil.
func isInvalidContext(c *gin.Context) bool {
	return c == nil || c.Request == nil
}

// isSkippedPath checks if the request URL matches any configured skip path.
func isSkippedPath(cfg *config.Config, c *gin.Context) bool {
	if c.Request.URL == nil {
		return false
	}
	return cfg.ShouldSkip(c.Request.URL.Path)
}

// shouldBypassTracing checks if request should bypass distributed tracing.
func shouldBypassTracing(cfg *config.Config, c *gin.Context) bool {
	if isTraceDisabled(cfg) || isInvalidContext(c) {
		return true
	}
	return isSkippedPath(cfg, c)
}

// canInjectHeaders checks if response headers can be written.
func canInjectHeaders(propagator propagation.TextMapPropagator, c *gin.Context) bool {
	if propagator == nil || isInvalidContext(c) {
		return false
	}
	return c.Writer != nil
}

// injectResponseHeaders injects trace context into outgoing response headers.
func injectResponseHeaders(propagator propagation.TextMapPropagator, c *gin.Context) {
	if !canInjectHeaders(propagator, c) {
		return
	}
	propagator.Inject(c.Request.Context(), propagation.HeaderCarrier(c.Writer.Header()))
}

// startTraceSpan begins a server span with initial attributes and returns updated context.
func startTraceSpan(
	tracer trace.Tracer,
	prop propagation.TextMapPropagator,
	c *gin.Context,
) (context.Context, trace.Span) {
	ctx := prop.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
	spanName := buildInitialSpanName(c)
	attrs := buildInitialAttributes(c)
	opts := []trace.SpanStartOption{trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attrs...)}
	return tracer.Start(ctx, spanName, opts...)
}

// traceHandler wraps HTTP execution inside an OpenTelemetry server span.
func traceHandler(tracer trace.Tracer, propagator propagation.TextMapPropagator, c *gin.Context) {
	ctx, span := startTraceSpan(tracer, propagator, c)
	defer func() { span.End() }()
	c.Request = c.Request.WithContext(ctx)
	injectTraceContext(c, span)
	injectResponseHeaders(propagator, c)
	c.Next()
	enrichSpanAfterNext(span, c)
}

// Middleware creates a Gin middleware for OpenTelemetry distributed tracing.
func Middleware(cfg *config.Config) gin.HandlerFunc {
	if isTraceDisabled(cfg) {
		return func(c *gin.Context) { c.Next() }
	}
	tracer := otel.GetTracerProvider().Tracer(DefaultTracerName)
	propagator := otel.GetTextMapPropagator()
	return func(c *gin.Context) {
		if shouldBypassTracing(cfg, c) {
			c.Next()
			return
		}
		traceHandler(tracer, propagator, c)
	}
}

// TraceMiddleware is an alias for Middleware to support both naming conventions.
func TraceMiddleware(cfg *config.Config) gin.HandlerFunc {
	return Middleware(cfg)
}
