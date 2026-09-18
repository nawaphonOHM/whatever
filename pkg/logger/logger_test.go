// Package logger provides unit tests for structured HTTP request logging.
package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testReqIDVal = "test-req-123"
	status200Val = 200
)

// setupLoggedEngine creates an engine with mock request ID and logger.
func setupLoggedEngine(buf *bytes.Buffer) *gin.Engine {
	logger := slog.New(slog.NewJSONHandler(buf, nil))
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(RequestIDKey, testReqIDVal)
		c.Writer.Header().Set(HeaderXRequestID, testReqIDVal)
		c.Next()
	})
	r.Use(WithLogger(logger))
	r.GET("/items/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

// verifyLoggedFields validates field content in parsed JSON log entry.
func verifyLoggedFields(t *testing.T, logEntry map[string]any) {
	assert.Equal(t, httpRequestMsg, logEntry["msg"])
	assert.Equal(t, "INFO", logEntry["level"])
	assert.Equal(t, "GET", logEntry["method"])
	assert.Equal(t, "/items/test", logEntry["path"])
	assert.Equal(t, "query=1", logEntry["query"])
	assert.Equal(t, float64(status200Val), logEntry["status"])
	assert.Equal(t, testReqIDVal, logEntry["request_id"])
}

// executeItemsRequest executes GET request against test route.
func executeItemsRequest(
	t *testing.T,
	r *gin.Engine,
) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/items/test?query=1", nil)
	require.NoError(t, err)
	req.Header.Set(HeaderXRequestID, testReqIDVal)
	r.ServeHTTP(w, req)
	return w
}

// TestLogger_StructuredLogging tests fields produced in structured logs.
func TestLogger_StructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggedEngine(&buf)

	w := executeItemsRequest(t, r)
	assert.Equal(t, http.StatusOK, w.Code)

	var logEntry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))
	verifyLoggedFields(t, logEntry)
}

// setupSkipPathEngine sets up an engine with skip paths configured.
func setupSkipPathEngine(buf *bytes.Buffer) *gin.Engine {
	logger := slog.New(slog.NewJSONHandler(buf, nil))
	r := gin.New()
	r.Use(WithConfig(Config{
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
