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

func TestLogger_StructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Simulate request-id middleware writing into the shared context key.
		c.Set(RequestIDKey, "test-req-123")
		c.Writer.Header().Set(HeaderXRequestID, "test-req-123")
		c.Next()
	})
	r.Use(WithLogger(logger))

	r.GET("/items/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/items/test?query=1", nil)
	req.Header.Set(HeaderXRequestID, "test-req-123")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err, "log output should be valid JSON: %s", buf.String())

	assert.Equal(t, "HTTP Request", logEntry["msg"])
	assert.Equal(t, "INFO", logEntry["level"])
	assert.Equal(t, "GET", logEntry["method"])
	assert.Equal(t, "/items/test", logEntry["path"])
	assert.Equal(t, "query=1", logEntry["query"])
	assert.Equal(t, float64(200), logEntry["status"])
	assert.Equal(t, "test-req-123", logEntry["request_id"])
}

func TestLogger_SkipPaths(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

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

	// Request skipped path
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, buf.String())

	// Request non-skipped path
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/other", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.NotEmpty(t, buf.String())
}

func TestLogger_ErrorLogLevels(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		expectedLevel string
	}{
		{"ClientError", http.StatusBadRequest, "WARN"},
		{"ServerError", http.StatusInternalServerError, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			r := gin.New()
			r.Use(WithLogger(logger))
			r.GET("/status", func(c *gin.Context) {
				c.Status(tt.status)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/status", nil)
			r.ServeHTTP(w, req)

			var logEntry map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &logEntry)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedLevel, logEntry["level"])
		})
	}
}
