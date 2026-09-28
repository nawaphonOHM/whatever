package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// findAttribute searches span attributes for a specific key.
func findAttribute(span sdktrace.ReadOnlySpan, key attribute.Key) (attribute.Value, bool) {
	for _, kv := range span.Attributes() {
		if kv.Key == key {
			return kv.Value, true
		}
	}
	return attribute.Value{}, false
}

func verifyRequestRouteAttributes(t *testing.T, span sdktrace.ReadOnlySpan) {
	methodVal, ok := findAttribute(span, AttrHTTPMethod)
	require.True(t, ok)
	assert.Equal(t, http.MethodPost, methodVal.AsString())

	routeVal, ok := findAttribute(span, AttrHTTPRoute)
	require.True(t, ok)
	assert.Equal(t, "/api/v1/orders/:id", routeVal.AsString())

	targetVal, ok := findAttribute(span, AttrHTTPTarget)
	require.True(t, ok)
	assert.Equal(t, "/api/v1/orders/12345?priority=high", targetVal.AsString())
}

func verifyStatusAndClientAttributes(t *testing.T, span sdktrace.ReadOnlySpan) {
	statusVal, ok := findAttribute(span, AttrHTTPStatusCode)
	require.True(t, ok)
	assert.Equal(t, int64(http.StatusCreated), statusVal.AsInt64())

	ipVal, ok := findAttribute(span, AttrClientAddress)
	require.True(t, ok)
	assert.Equal(t, "192.0.2.1", ipVal.AsString())
}

func buildStandardOrderRouter() *gin.Engine {
	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	router.POST("/api/v1/orders/:id", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})
	return router
}

func executeStandardOrderRequest(router *gin.Engine) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/12345?priority=high", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func verifyStandardSpan(t *testing.T, span sdktrace.ReadOnlySpan) {
	assert.Equal(t, "POST /api/v1/orders/:id", span.Name())
	assert.Equal(t, codes.Ok, span.Status().Code)
	verifyRequestRouteAttributes(t, span)
	verifyStatusAndClientAttributes(t, span)
}

func TestMiddleware_Attributes_Standard(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	router := buildStandardOrderRouter()
	rec := executeStandardOrderRequest(router)

	assert.Equal(t, http.StatusCreated, rec.Code)
	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	verifyStandardSpan(t, spans[0])
}
