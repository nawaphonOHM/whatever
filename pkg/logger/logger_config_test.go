// Package logger provides unit tests for structured HTTP request logging configuration.
package logger

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupSkipPathEngine sets up an engine with skip paths configured.
func setupSkipPathEngine(buf *bytes.Buffer) *gin.Engine {
	logger := slog.New(slog.NewJSONHandler(buf, nil))
	r := gin.New()
	r.Use(WithConfig(&Config{
		Logger:    logger,
		SkipPaths: []string{"/health"},
	}))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})
	r.GET("/other", func(c *gin.Context) {
		c.String(http.StatusOK, "other")
	})
	return r
}

// executeTestRequest runs a GET request against router and returns recorder.
func executeTestRequest(
	t *testing.T,
	r *gin.Engine,
	path string,
) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)
	return w
}

// TestLogger_SkipPaths tests skipping configured paths from logging.
func TestLogger_SkipPaths(t *testing.T) {
	var buf bytes.Buffer
	r := setupSkipPathEngine(&buf)

	// Execute skipped request and verify no log output
	w := executeTestRequest(t, r, "/health")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, buf.String())

	// Execute non-skipped request and verify log output
	w2 := executeTestRequest(t, r, "/other")
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.NotEmpty(t, buf.String())
}

// TestLogger_NilConfig verifies that passing nil config defaults gracefully.
func TestLogger_NilConfig(t *testing.T) {
	r := gin.New()
	r.Use(WithConfig(nil))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := executeTestRequest(t, r, "/test")
	assert.Equal(t, http.StatusOK, w.Code)
}
