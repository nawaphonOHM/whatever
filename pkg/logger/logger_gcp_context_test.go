package logger_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logger"
)

const (
	testContextReqID   = "req-xyz-999"
	testContextTraceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	testContextSpanID  = "00f067aa0ba902b7"
	testContextPath    = "/api/v1/resource/123"
	testContextAgent   = "custom-runner"
	testContextProto   = "HTTP/2.0"
	testErrContextMsg  = "context test error"
)

func setupContextEngine(buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	gcpLogger := slog.New(core.NewGCPHandler(buf, nil))
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(logger.RequestIDKey, testContextReqID)
		c.Set("TraceID", testContextTraceID)
		c.Set("SpanID", testContextSpanID)
		if err := c.Error(errors.New(testErrContextMsg)); err != nil {
			c.Status(http.StatusInternalServerError)
		}
		c.Next()
	})
	r.Use(logger.WithLogger(gcpLogger))
	r.DELETE(testContextPath, func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	return r
}

func verifyContextTopLevel(t *testing.T, entry map[string]any) {
	assert.Equal(t, "INFO", entry["severity"])
	assert.Equal(t, testContextReqID, entry["request_id"])
	assert.Equal(t, testContextTraceID, entry["trace_id"])
	assert.Equal(t, testContextSpanID, entry["span_id"])
	assert.Contains(t, entry["errors"], testErrContextMsg)
}

func verifyContextHTTPRequest(t *testing.T, httpReq map[string]any) {
	assert.Equal(t, "DELETE", httpReq["requestMethod"])
	assert.Equal(t, testContextPath, httpReq["requestUrl"])
	assert.Equal(t, float64(http.StatusNoContent), httpReq["status"])
	assert.Equal(t, testContextProto, httpReq["protocol"])
	assert.Equal(t, testContextAgent, httpReq["userAgent"])
}

func executeContextRequest(t *testing.T, r *gin.Engine, buf *bytes.Buffer) map[string]any {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, testContextPath, nil)
	req.Proto = testContextProto
	req.Header.Set("User-Agent", testContextAgent)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	return entry
}

func verifyContextExclusions(t *testing.T, entry map[string]any) {
	flatKeys := []string{"method", "path", "query", "status", "latency", "latency_ms", "bytes_out", "protocol"}
	for _, key := range flatKeys {
		assert.Nil(t, entry[key], "flat key %s should not be present in root payload", key)
	}
}

func TestLogger_GCP_EndToEnd_WithTraceAndErrors(t *testing.T) {
	var buf bytes.Buffer
	r := setupContextEngine(&buf)
	entry := executeContextRequest(t, r, &buf)

	verifyContextTopLevel(t, entry)
	httpReq, ok := entry["httpRequest"].(map[string]any)
	require.True(t, ok)
	verifyContextHTTPRequest(t, httpReq)
	verifyContextExclusions(t, entry)
}
