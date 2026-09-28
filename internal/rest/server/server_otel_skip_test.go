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

	assert.Equal(t, http.StatusOK, performGet(srv, "/v1/items/skip").Code)
	assert.Empty(t, exporter.GetSpans().Snapshots())

	assert.Equal(t, http.StatusOK, performGet(srv, "/v1/items/real").Code)
	spans := exporter.GetSpans().Snapshots()
	require.Len(t, spans, 1)
	assert.Equal(t, "GET /v1/items/:id", spans[0].Name())
}

func buildCustomSkipServer(t *testing.T) (*tracetest.InMemoryExporter, func(), *Server) {
	t.Helper()
	exporter, cleanup := setupTestTracer(t)
	cfg := DefaultConfig()
	cfg.OTel.SkipPaths = []string{"/v1/items/skip"}
	srv := buildTracedServer(t, cfg, itemsBluePrint())
	return exporter, cleanup, srv
}
