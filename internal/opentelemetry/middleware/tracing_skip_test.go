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
)

func buildSkipTestRouter(cfg *config.Config) *gin.Engine {
	router := setupTestEngine(Middleware(cfg))
	handler := func(c *gin.Context) { c.Status(http.StatusOK) }
	router.GET("/health", handler)
	router.GET("/ready", handler)
	router.GET("/api/v1/data", handler)
	return router
}

func verifySkippedEndpoint(t *testing.T, router *gin.Engine, exporter *tracetest.InMemoryExporter, path string) {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, getRecordedSpans(t, exporter))
}

func verifyActiveEndpoint(
	t *testing.T,
	router *gin.Engine,
	exporter *tracetest.InMemoryExporter,
) sdktrace.ReadOnlySpan {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/data", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	return spans[0]
}

func TestMiddleware_SkipPaths(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	cfg := &config.Config{
		SkipPaths: []string{"/health", "/ready"},
		Enabled:   true,
	}

	router := buildSkipTestRouter(cfg)
	verifySkippedEndpoint(t, router, exporter, "/health")
	verifySkippedEndpoint(t, router, exporter, "/ready")

	span := verifyActiveEndpoint(t, router, exporter)
	assert.Equal(t, "GET /api/v1/data", span.Name())
}
