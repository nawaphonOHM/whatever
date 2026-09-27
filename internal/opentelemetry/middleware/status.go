package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// recordSingleError records non-nil gin error on the span.
func recordSingleError(span trace.Span, ginErr *gin.Error) {
	if ginErr == nil {
		return
	}
	if ginErr.Err != nil {
		span.RecordError(ginErr.Err)
	}
}

// iterateErrors records all errors from slice onto the span.
func iterateErrors(span trace.Span, errs []*gin.Error) {
	for _, ginErr := range errs {
		recordSingleError(span, ginErr)
	}
}

// recordSpanErrors records any Gin context errors on the active span.
func recordSpanErrors(span trace.Span, c *gin.Context) {
	if span == nil {
		return
	}
	if c == nil {
		return
	}
	iterateErrors(span, c.Errors)
}

// errorDescription builds a descriptive error message for 5xx status codes.
func errorDescription(c *gin.Context, statusCode int) string {
	if c != nil && len(c.Errors) > 0 {
		return c.Errors.String()
	}
	return http.StatusText(statusCode)
}

// recordSpanStatus sets the span status according to the HTTP response code.
func recordSpanStatus(span trace.Span, c *gin.Context, statusCode int) {
	if span == nil {
		return
	}
	recordSpanErrors(span, c)
	if statusCode >= statusServerError {
		span.SetStatus(codes.Error, errorDescription(c, statusCode))
		return
	}
	span.SetStatus(codes.Ok, "")
}
