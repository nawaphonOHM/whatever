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

func setupNoContentRouter(buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	gcpLogger := slog.New(core.NewGCPHandler(buf, nil))
	r := gin.New()
	r.Use(logger.WithLogger(gcpLogger))
	r.POST("/empty-post", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	return r
}

func verifyNoContentReq(t *testing.T, httpReq map[string]any) {
	assert.Equal(t, "POST", httpReq["requestMethod"])
	assert.Equal(t, "/empty-post", httpReq["requestUrl"])
	assert.Equal(t, float64(http.StatusNoContent), httpReq["status"])
	assert.Nil(t, httpReq["responseSize"], "responseSize must be omitted for 204 No Content")
	latency, ok := httpReq["latency"].(string)
	require.True(t, ok)
	assert.True(t, strings.HasSuffix(latency, "s"))
}

func verifyNoContentPayload(t *testing.T, entry map[string]any) {
	assert.Equal(t, "INFO", entry["severity"])
	httpReq, ok := entry["httpRequest"].(map[string]any)
	require.True(t, ok)
	verifyNoContentReq(t, httpReq)

	flatKeys := []string{"method", "path", "query", "status", "latency", "latency_ms", "bytes_out", "protocol"}
	for _, key := range flatKeys {
		assert.Nil(t, entry[key], "flat key %s should not be present in root payload", key)
	}
}

func TestLogger_GCP_PostNoContent(t *testing.T) {
	var buf bytes.Buffer
	r := setupNoContentRouter(&buf)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/empty-post", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	verifyNoContentPayload(t, entry)
}
