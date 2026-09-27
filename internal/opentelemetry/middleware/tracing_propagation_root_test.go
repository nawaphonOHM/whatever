package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func verifyRootSpan(t *testing.T, span sdktrace.ReadOnlySpan, rec *httptest.ResponseRecorder) {
	assert.False(t, span.Parent().IsValid())
	assert.Contains(t, rec.Header().Get(HeaderTraceParent), span.SpanContext().TraceID().String())
}

func verifyNewRootResults(
	t *testing.T,
	handlerCtx context.Context,
	rec *httptest.ResponseRecorder,
	exporter *tracetest.InMemoryExporter,
) {
	assert.Equal(t, http.StatusOK, rec.Code)
	require.True(t, trace.SpanFromContext(handlerCtx).SpanContext().IsValid())
	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	verifyRootSpan(t, spans[0], rec)
}

func TestMiddleware_TraceContextPropagation_NewRoot(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	var handlerCtx context.Context
	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	router.GET("/root", func(c *gin.Context) {
		handlerCtx = c.Request.Context()
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/root", nil))
	verifyNewRootResults(t, handlerCtx, rec, exporter)
}

func executeInvalidHeaderRequest(router *gin.Engine) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/fallback", nil)
	req.Header.Set(HeaderTraceParent, "invalid-malformed-traceparent")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestMiddleware_TraceContextPropagation_InvalidHeader(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	router.GET("/fallback", func(c *gin.Context) { c.Status(http.StatusOK) })

	rec := executeInvalidHeaderRequest(router)
	assert.Equal(t, http.StatusOK, rec.Code)

	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	assert.False(t, spans[0].Parent().IsValid())
	assert.True(t, spans[0].SpanContext().IsValid())
}
