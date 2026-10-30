package logger_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

func setupLevelsRouter(buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	gcpLogger := slog.New(core.NewGCPHandler(buf, nil))
	r := gin.New()
	r.Use(logger.WithLogger(gcpLogger))
	r.GET("/bad-request", func(c *gin.Context) {
		c.String(http.StatusBadRequest, "invalid parameter")
	})
	r.GET("/server-error", func(c *gin.Context) {
		if err := c.Error(errors.New("db connection failure")); err != nil {
			c.String(http.StatusInternalServerError, "internal server error")
		}
	})
	return r
}

func verifyLevelExclusions(t *testing.T, entry map[string]any) {
	flatKeys := []string{"method", "path", "query", "status", "latency", "latency_ms", "bytes_out", "protocol"}
	for _, key := range flatKeys {
		assert.Nil(t, entry[key], "flat key %s should not be present in root payload", key)
	}
}

func verifyLevelHTTPReq(t *testing.T, entry map[string]any, sev string, status int) {
	assert.Equal(t, sev, entry["severity"])
	httpReq, ok := entry["httpRequest"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(status), httpReq["status"])
	assert.Equal(t, "GET", httpReq["requestMethod"])
	assert.Greater(t, httpReq["responseSize"], float64(0))

	latency, ok := httpReq["latency"].(string)
	require.True(t, ok)
	assert.True(t, strings.HasSuffix(latency, "s"))
	verifyLevelExclusions(t, entry)
}

func TestLogger_GCP_BadRequest_Level(t *testing.T) {
	var buf bytes.Buffer
	r := setupLevelsRouter(&buf)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bad-request", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	verifyLevelHTTPReq(t, entry, "WARNING", http.StatusBadRequest)
}

func TestLogger_GCP_ServerError_Level(t *testing.T) {
	var buf bytes.Buffer
	r := setupLevelsRouter(&buf)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/server-error", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	verifyLevelHTTPReq(t, entry, "ERROR", http.StatusInternalServerError)
	assert.Contains(t, entry["errors"], "db connection failure")
}
