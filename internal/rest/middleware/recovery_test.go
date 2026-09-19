// Package middleware provides unit tests for panic recovery middleware.
package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	panicMsg   = "something went critically wrong"
	panicPath  = "/panic"
	panicReqID = "panic-req-id"
)

// setupPanicRouter creates a test router with recovery middleware that panics.
func setupPanicRouter(logger *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(RequestID())
	r.Use(RecoveryWithLogger(logger))
	r.GET(panicPath, func(*gin.Context) {
		panic(panicMsg)
	})
	return r
}

// verifyProblemDetails checks RFC 9457 fields on unmarshaled response.
func verifyProblemDetails(t *testing.T, prob rest.ProblemDetails) {
	assert.Equal(t, http.StatusInternalServerError, prob.Status)
	assert.Equal(t, "about:blank", prob.Type)
	assert.Equal(t, "Internal Server Error", prob.Title)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", prob.Code)
	assert.Equal(t, "An unexpected internal server error occurred", prob.Detail)
	assert.Equal(t, panicPath, prob.Instance)
}

// verifyPanicResponse validates status code and RFC problem details.
func verifyPanicResponse(t *testing.T, w *httptest.ResponseRecorder) {
	// Verify HTTP status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Verify response content type
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

	var prob rest.ProblemDetails
	err := json.Unmarshal(w.Body.Bytes(), &prob)
	require.NoError(t, err)

	verifyProblemDetails(t, prob)
}

// verifyPanicLogs validates that the buffer contains required log entries.
func verifyPanicLogs(t *testing.T, buf *bytes.Buffer) {
	// Extract output from test buffer
	logOutput := buf.String()
	// Assert presence of recovery message
	assert.Contains(t, logOutput, "panic recovered during request processing")
	// Assert presence of panic detail
	assert.Contains(t, logOutput, panicMsg)
	// Assert presence of request correlation ID
	assert.Contains(t, logOutput, panicReqID)
}

// TestRecovery_HandlesPanic verifies that panics are caught properly.
func TestRecovery_HandlesPanic(t *testing.T) {
	var buf bytes.Buffer
	// Initialize logger capturing output to buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	r := setupPanicRouter(logger)

	// Create request with custom request ID
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, panicPath, nil)
	require.NoError(t, err)
	req.Header.Set(HeaderXRequestID, panicReqID)

	// Execute HTTP request ensuring no crash
	assert.NotPanics(t, func() {
		r.ServeHTTP(w, req)
	})

	// Verify HTTP response and logs
	verifyPanicResponse(t, w)
	verifyPanicLogs(t, &buf)
}
