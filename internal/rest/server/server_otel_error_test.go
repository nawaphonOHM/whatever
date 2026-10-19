package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// performGet sends a GET request to the given path and returns recorder.
func performGet(s *Server, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

// TestServer_OTel_NotFound verifies spans for 404 RFC9457 responses.
func TestServer_OTel_NotFound(t *testing.T) {
	assertErrorSpan(t, http.MethodGet, "/api/v1/nonexistent", http.StatusNotFound)
}

// TestServer_OTel_MethodNotAllowed verifies spans for 405 RFC9457 responses.
func TestServer_OTel_MethodNotAllowed(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	srv := buildTracedServer(t, DefaultConfig(), itemsBluePrint())
	w405 := httptest.NewRecorder()
	srv.Engine.ServeHTTP(w405, httptest.NewRequest(http.MethodPost, "/api/v1/items/1", nil))
	assert.Equal(t, http.StatusMethodNotAllowed, w405.Code)

	spans := exporter.GetSpans().Snapshots()
	require.Len(t, spans, 1)
	assert.Equal(t, "POST", spans[0].Name())
}

func assertErrorSpan(t *testing.T, method, path string, status int) {
	t.Helper()
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	srv := buildTracedServer(t, DefaultConfig(), itemsBluePrint())
	w := httptest.NewRecorder()
	srv.Engine.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	assert.Equal(t, status, w.Code)
	assertRecordedErrorSpan(t, exporter, method)
}

func assertRecordedErrorSpan(
	t *testing.T,
	exporter *tracetest.InMemoryExporter,
	method string,
) {
	t.Helper()
	spans := exporter.GetSpans().Snapshots()
	require.Len(t, spans, 1)
	assert.Equal(t, method, spans[0].Name())
}

// TestServer_OTel_DisabledTelemetry verifies disabled telemetry creates no spans.
func TestServer_OTel_DisabledTelemetry(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	cfg := DefaultConfig()
	cfg.OTel.Enabled = false
	srv := buildTracedServer(t, cfg, itemsBluePrint())

	w := performGet(srv, "/api/v1/items/42")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, exporter.GetSpans().Snapshots())
}

// TestServer_OTel_HealthSkipPaths verifies that /health and /ready do not generate spans.
func TestServer_OTel_HealthSkipPaths(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	srv := buildTracedServer(t, DefaultConfig(), itemsBluePrint())
	assert.Equal(t, http.StatusOK, performGet(srv, ReservedHealthPath).Code)
	assert.Equal(t, http.StatusOK, performGet(srv, ReservedReadyPath).Code)
	assert.Empty(t, exporter.GetSpans().Snapshots())
}

// TestServer_OTel_CustomSkipPaths verifies custom skip paths bypass span creation.
func TestServer_OTel_CustomSkipPaths(t *testing.T) {
	exporter, cleanup, srv := buildCustomSkipServer(t)
	defer cleanup()

	assert.Equal(t, http.StatusOK, performGet(srv, "/api/v1/items/skip").Code)
	assert.Empty(t, exporter.GetSpans().Snapshots())

	assert.Equal(t, http.StatusOK, performGet(srv, "/api/v1/items/real").Code)
	spans := exporter.GetSpans().Snapshots()
	require.Len(t, spans, 1)
	assert.Equal(t, "GET /api/v1/items/:id", spans[0].Name())
}

func buildCustomSkipServer(t *testing.T) (*tracetest.InMemoryExporter, func(), *Server) {
	t.Helper()
	exporter, cleanup := setupTestTracer(t)
	cfg := DefaultConfig()
	cfg.OTel.SkipPaths = []string{"/api/v1/items/skip"}
	srv := buildTracedServer(t, cfg, itemsBluePrint())
	return exporter, cleanup, srv
}
