package middleware

import (
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

const (
	testTraceParent  = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	expectedTraceID  = "4bf92f3577b34da6a3ce929d0e0e4736"
	expectedParentID = "00f067aa0ba902b7"
)

type propagationResult struct {
	traceID        string
	spanID         string
	handlerHasSpan bool
}

func recordPropagationHandler(res *propagationResult) gin.HandlerFunc {
	return func(c *gin.Context) {
		res.traceID = GetTraceID(c)
		res.spanID = GetSpanID(c)
		res.handlerHasSpan = trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid()
		c.Status(http.StatusOK)
	}
}

func sendParentPropagationRequest(router *gin.Engine) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	req.Header.Set(HeaderTraceParent, testTraceParent)
	req.Header.Set(HeaderTraceState, "congo=t61rcWkgMzE")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func verifyRecordedParentSpan(t *testing.T, span sdktrace.ReadOnlySpan) {
	assert.Equal(t, expectedTraceID, span.SpanContext().TraceID().String())
	assert.Equal(t, expectedParentID, span.Parent().SpanID().String())
	assert.Equal(t, trace.SpanKindServer, span.SpanKind())
}

func verifyParentTestResult(
	t *testing.T,
	res *propagationResult,
	rec *httptest.ResponseRecorder,
	exporter *tracetest.InMemoryExporter,
) {
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, res.handlerHasSpan)
	assert.Equal(t, expectedTraceID, res.traceID)
	assert.NotEmpty(t, res.spanID)
	assert.Contains(t, rec.Header().Get(HeaderTraceParent), expectedTraceID)

	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	verifyRecordedParentSpan(t, spans[0])
}

func TestMiddleware_TraceContextPropagation_WithParent(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	var res propagationResult
	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	router.GET("/users/:id", recordPropagationHandler(&res))

	rec := sendParentPropagationRequest(router)
	verifyParentTestResult(t, &res, rec, exporter)
}
