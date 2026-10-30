package logger_test

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

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logger"
)

func setupEncodedPathEngine(buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	gcpLogger := slog.New(core.NewGCPHandler(buf, nil))
	r := gin.New()
	r.Use(logger.WithLogger(gcpLogger))
	r.GET("/files/*filepath", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"file": c.Param("filepath")})
	})
	return r
}

func verifyEncodedHTTPReq(t *testing.T, httpReq map[string]any, targetURI string) {
	assert.Equal(t, "GET", httpReq["requestMethod"])
	assert.Equal(t, targetURI, httpReq["requestUrl"])
	assert.Equal(t, float64(http.StatusOK), httpReq["status"])
}

func verifyEncodedExclusions(t *testing.T, entry map[string]any) {
	assert.Equal(t, "INFO", entry["severity"])
	flatKeys := []string{"method", "path", "query", "status", "latency", "latency_ms", "bytes_out"}
	for _, key := range flatKeys {
		assert.Nil(t, entry[key], "flat key %s should not be present in root payload", key)
	}
}

func executeEncodedReq(t *testing.T, r *gin.Engine, uri string, buf *bytes.Buffer) map[string]any {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	return entry
}

func TestLogger_GCP_EncodedPathSegment(t *testing.T) {
	var buf bytes.Buffer
	r := setupEncodedPathEngine(&buf)

	targetURI := "/files/hello%20world%2Fsubfolder?tag=go%2B%2B&order=asc"
	entry := executeEncodedReq(t, r, targetURI, &buf)
	httpReq, ok := entry["httpRequest"].(map[string]any)
	require.True(t, ok)

	verifyEncodedHTTPReq(t, httpReq, targetURI)
	verifyEncodedExclusions(t, entry)
}
