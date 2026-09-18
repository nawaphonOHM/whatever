// Package logger provides unit tests for log level classification.
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

// setupStatusEngine creates router configured with status endpoint.
func setupStatusEngine(buf *bytes.Buffer, status int) *gin.Engine {
	logger := slog.New(slog.NewJSONHandler(buf, nil))
	r := gin.New()
	r.Use(WithLogger(logger))
	r.GET("/status", func(c *gin.Context) {
		c.Status(status)
	})
	return r
}

// runErrorLogLevelCase runs a subtest checking log level for a given status.
func runErrorLogLevelCase(t *testing.T, expectedLevel string, status int) {
	var buf bytes.Buffer
	r := setupStatusEngine(&buf, status)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)

	var logEntry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))
	assert.Equal(t, expectedLevel, logEntry["level"])
}

// TestLogger_ErrorLogLevels tests level mapping based on status code.
func TestLogger_ErrorLogLevels(t *testing.T) {
	t.Run("ClientError", func(t *testing.T) {
		runErrorLogLevelCase(t, "WARN", http.StatusBadRequest)
	})
	t.Run("ServerError", func(t *testing.T) {
		runErrorLogLevelCase(t, "ERROR", http.StatusInternalServerError)
	})
}
