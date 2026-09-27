package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func buildTestSpanContext(t *testing.T) context.Context {
	traceID, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("00f067aa0ba902b7")
	require.NoError(t, err)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(context.Background(), spanCtx)
}

func TestContext_GetTraceID_GetSpanID_NilSafety(t *testing.T) {
	assert.Empty(t, GetTraceID(nil))
	assert.Empty(t, GetSpanID(nil))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	assert.Empty(t, GetTraceID(c))
	assert.Empty(t, GetSpanID(c))
}

func TestContext_GetTraceID_GetSpanID_FromGinKeys(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(TraceIDKey, "trace-12345")
	c.Set(SpanIDKey, "span-67890")

	assert.Equal(t, "trace-12345", GetTraceID(c))
	assert.Equal(t, "span-67890", GetSpanID(c))
}

func TestContext_GetTraceID_GetSpanID_FromSpanContext(t *testing.T) {
	ctx := buildTestSpanContext(t)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
	require.NoError(t, err)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", GetTraceID(c))
	assert.Equal(t, "00f067aa0ba902b7", GetSpanID(c))
}

func TestContext_NilHelpers(t *testing.T) {
	injectTraceContext(nil, nil)
	enrichSpanAfterNext(nil, nil)
	recordSpanErrors(nil, nil)
	recordSpanStatus(nil, nil, http.StatusOK)

	assert.Equal(t, "HTTP", buildInitialSpanName(nil))
	assert.Empty(t, getTargetURI(nil))
	assert.Nil(t, buildInitialAttributes(nil))
}
