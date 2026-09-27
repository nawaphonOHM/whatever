package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func verifyUnmatchedSpan(t *testing.T, span sdktrace.ReadOnlySpan) {
	assert.Equal(t, "GET", span.Name())
	assert.Equal(t, codes.Ok, span.Status().Code)
	val, ok := findAttribute(span, AttrHTTPStatusCode)
	require.True(t, ok)
	assert.Equal(t, int64(http.StatusNotFound), val.AsInt64())
}

func TestMiddleware_Attributes_UnmatchedRoute(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/unmatched/path", nil))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	verifyUnmatchedSpan(t, spans[0])
}
