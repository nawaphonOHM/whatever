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

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logger"
)

const testAPIRoute = "/api/test"

func createConfigEngine(cfg *config.Config, buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := core.BuildHandler(cfg, buf, nil)
	customLogger := slog.New(handler)
	r := gin.New()
	r.Use(logger.WithLogger(customLogger))
	r.GET(testAPIRoute, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "value"})
	})
	r.POST("/api/empty", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	return r
}

func executeConfigReq(t *testing.T, r *gin.Engine, req *http.Request, buf *bytes.Buffer) map[string]any {
	t.Helper()
	buf.Reset()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	return entry
}

func verifyGCPEnvDetection(t *testing.T, entry map[string]any) {
	assert.Equal(t, "INFO", entry["severity"])
	httpReq, ok := entry["httpRequest"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "GET", httpReq["requestMethod"])
	assert.Equal(t, testAPIRoute, httpReq["requestUrl"])
	assert.Nil(t, entry["method"])
}

func TestLogger_GCP_BuildHandler_EnvDetection(t *testing.T) {
	var buf bytes.Buffer
	cfg := &config.Config{
		CloudProvider:           "gcp",
		DisableTraceCorrelation: true,
		Format:                  "json",
		Output:                  "stdout",
	}
	r := createConfigEngine(cfg, &buf)
	req := httptest.NewRequest(http.MethodGet, testAPIRoute, nil)
	entry := executeConfigReq(t, r, req, &buf)
	verifyGCPEnvDetection(t, entry)
}

func verifyNonGCPAbsentFields(t *testing.T, entry map[string]any) {
	assert.Nil(t, entry["protocol"], "protocol should not be present in non-GCP format")
	assert.Nil(t, entry["request_url"], "request_url should not be present in non-GCP format")
	assert.Nil(t, entry["httpRequest"], "httpRequest should not be present in non-GCP format")
}

func verifyStandardNonGCPEntry(t *testing.T, entry map[string]any) {
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, testAPIRoute, entry["path"])
	assert.Equal(t, "", entry["query"])
	assert.Equal(t, float64(http.StatusOK), entry["status"])
	assert.NotNil(t, entry["latency_ms"])
	assert.NotNil(t, entry["latency"])
	assert.Greater(t, entry["bytes_out"], float64(0))
	verifyNonGCPAbsentFields(t, entry)
}

func TestLogger_NonGCP_PreservesStandardFormat(t *testing.T) {
	var buf bytes.Buffer
	cfg := &config.Config{
		CloudProvider:           "",
		DisableTraceCorrelation: true,
		Format:                  "json",
		Output:                  "stdout",
	}
	r := createConfigEngine(cfg, &buf)
	reqGet := httptest.NewRequest(http.MethodGet, testAPIRoute, nil)
	entry := executeConfigReq(t, r, reqGet, &buf)
	verifyStandardNonGCPEntry(t, entry)

	reqPost := httptest.NewRequest(http.MethodPost, "/api/empty", nil)
	entryNoContent := executeConfigReq(t, r, reqPost, &buf)
	assert.Equal(t, float64(-1), entryNoContent["bytes_out"], "non-GCP 204 must preserve raw -1 bytes_out")
	assert.Nil(t, entryNoContent["protocol"])
}
