// Package health provides unit tests for internal health probe endpoints.
package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testVersion defines the test version string used across health tests.
const testVersion = "2.0.0"

// init configures gin to test mode for all test executions.
func init() {
	gin.SetMode(gin.TestMode)
}

// performProbeRequest executes a GET request against the given path.
func performProbeRequest(h *Handler, path string) *httptest.ResponseRecorder {
	r := gin.New()
	if path == "/health" {
		r.GET("/health", h.Health)
	} else {
		r.GET("/ready", h.Ready)
	}
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// testResponseEnvelope represents the expected envelope structure in tests.
type testResponseEnvelope[T any] struct {
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

// verifyProbeResponse validates the HTTP status and JSON response body.
func verifyProbeResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus string,
) {
	assert.Equal(t, http.StatusOK, w.Code)
	var resp testResponseEnvelope[Status]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.Equal(t, expectedStatus, resp.Data.Status)
	assert.Equal(t, testVersion, resp.Data.Version)
	assert.False(t, resp.Data.Timestamp.IsZero())
}

// TestHealthHandler_Health tests the liveness probe endpoint.
func TestHealthHandler_Health(t *testing.T) {
	h := New(testVersion)
	w := performProbeRequest(h, "/health")
	verifyProbeResponse(t, w, "up")
}

// TestHealthHandler_Ready tests the readiness probe endpoint.
func TestHealthHandler_Ready(t *testing.T) {
	h := New(testVersion)
	w := performProbeRequest(h, "/ready")
	verifyProbeResponse(t, w, "ready")
}
