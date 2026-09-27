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
	"go.opentelemetry.io/otel/trace"
)

const (
	testTraceIDHex = "4bf92f3577b34da6a3ce929d0e0e4736"
	testSpanIDHex  = "00f067aa0ba902b7"
)

// setupMockTraceContext attaches mock trace and span IDs to the request context.
func setupMockTraceContext(t *testing.T, req *http.Request) *http.Request {
	t.Helper()
	tid, err := trace.TraceIDFromHex(testTraceIDHex)
	require.NoError(t, err)
	sid, err := trace.SpanIDFromHex(testSpanIDHex)
	require.NoError(t, err)

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(req.Context(), sc)
	return req.WithContext(ctx)
}

// setupTraceEngine creates a test engine for trace correlation testing.
func setupTraceEngine(logger *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(WithLogger(logger))
	r.GET("/traced", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

// setupKeysEngine creates a test engine with Gin context keys.
func setupKeysEngine(logger *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(traceIDKey, testTraceIDHex)
		c.Set(spanIDKey, testSpanIDHex)
		c.Next()
	})
	r.Use(WithLogger(logger))
	r.GET("/keys", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

// verifyLogTrace asserts trace_id and span_id are present in parsed log JSON.
func verifyLogTrace(t *testing.T, data []byte) {
	t.Helper()
	var logEntry map[string]any
	require.NoError(t, json.Unmarshal(data, &logEntry))
	assert.Equal(t, testTraceIDHex, logEntry["trace_id"])
	assert.Equal(t, testSpanIDHex, logEntry["span_id"])
}

// TestLogger_TraceCorrelation verifies trace_id and span_id are included in logs when trace context is present.
func TestLogger_TraceCorrelation(t *testing.T) {
	var buf bytes.Buffer
	r := setupTraceEngine(slog.New(slog.NewJSONHandler(&buf, nil)))
	req := setupMockTraceContext(t, httptest.NewRequest(http.MethodGet, "/traced", nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	verifyLogTrace(t, buf.Bytes())
}

// TestLogger_GinKeyTraceCorrelation verifies trace_id and span_id from Gin context keys.
func TestLogger_GinKeyTraceCorrelation(t *testing.T) {
	var buf bytes.Buffer
	r := setupKeysEngine(slog.New(slog.NewJSONHandler(&buf, nil)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/keys", nil))

	assert.Equal(t, http.StatusOK, w.Code)
	verifyLogTrace(t, buf.Bytes())
}

// TestLogger_NilContextSafety tests nil safety of trace extraction functions.
func TestLogger_NilContextSafety(t *testing.T) {
	assert.Empty(t, extractTraceID(nil))
	assert.Empty(t, extractSpanID(nil))
	assert.Empty(t, extractSpanTraceID(nil))
	assert.Empty(t, extractSpanContextID(nil))
	assert.Empty(t, extractKey(nil, traceIDKey))
	assert.Empty(t, traceAttrs(nil))
}
