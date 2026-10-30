package logger_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logger"
)

const (
	testGCPReqID     = "test-gcp-req-456"
	testGCPUserAgent = "test-agent-v1"
	testGCPClientIP  = "203.0.113.195"
	testGCPURL       = "/api/v1/search?query=hello%20world&tag=a%26b"
)

func setupGCPEngine(buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	gcpLogger := slog.New(core.NewGCPHandler(buf, nil))
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(logger.RequestIDKey, testGCPReqID)
		c.Writer.Header().Set(logger.HeaderXRequestID, testGCPReqID)
		c.Next()
	})
	r.Use(logger.WithLogger(gcpLogger))
	r.GET("/api/v1/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func verifyGCPHTTPRequestFields(t *testing.T, httpReq map[string]any) {
	assert.Equal(t, "GET", httpReq["requestMethod"])
	assert.Equal(t, testGCPURL, httpReq["requestUrl"])
	assert.Equal(t, float64(http.StatusOK), httpReq["status"])
	assert.Equal(t, testGCPUserAgent, httpReq["userAgent"])
	assert.Equal(t, testGCPClientIP, httpReq["remoteIp"])
	assert.Equal(t, "HTTP/1.1", httpReq["protocol"])
	assert.Greater(t, httpReq["responseSize"], float64(0))
	latency, ok := httpReq["latency"].(string)
	require.True(t, ok)
	assert.True(t, strings.HasSuffix(latency, "s"))
}

func verifyGCPPayloadExclusions(t *testing.T, entry map[string]any) {
	flatKeys := []string{
		"method", "path", "query", "status", "latency",
		"latency_ms", "client_ip", "user_agent", "bytes_out", "protocol",
	}
	for _, key := range flatKeys {
		assert.Nil(t, entry[key], "key %s should not be present in root payload", key)
	}
	assert.Equal(t, testGCPReqID, entry["request_id"])
	assert.Equal(t, "INFO", entry["severity"])
	assert.Equal(t, "HTTP Request", entry["message"])
	assert.NotEmpty(t, entry["timestamp"])
}

func sendGCPTestRequest(r *gin.Engine) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, testGCPURL, nil)
	req.Header.Set("User-Agent", testGCPUserAgent)
	req.Header.Set("X-Forwarded-For", testGCPClientIP)
	r.ServeHTTP(w, req)
	return w, req
}

func TestLogger_GCP_EndToEnd_RequestPayload(t *testing.T) {
	var buf bytes.Buffer
	r := setupGCPEngine(&buf)

	w, _ := sendGCPTestRequest(r)
	assert.Equal(t, http.StatusOK, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	httpReq, ok := entry["httpRequest"].(map[string]any)
	require.True(t, ok)
	verifyGCPHTTPRequestFields(t, httpReq)
	verifyGCPPayloadExclusions(t, entry)
}
