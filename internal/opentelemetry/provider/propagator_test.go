package provider

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestCompositePropagator_Fields(t *testing.T) {
	prop := NewCompositePropagator()
	require.NotNil(t, prop)

	fields := prop.Fields()
	assert.Contains(t, fields, "traceparent")
	assert.Contains(t, fields, "tracestate")
	assert.Contains(t, fields, "baggage")
}

func createTestSpanContext(t *testing.T) (context.Context, trace.TraceID, trace.SpanID) {
	traceID, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("00f067aa0ba902b7")
	require.NoError(t, err)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(context.Background(), spanCtx), traceID, spanID
}

func assertSpanContextMatches(
	t *testing.T,
	extractedCtx context.Context,
	traceID trace.TraceID,
	spanID trace.SpanID,
) {
	extractedSpanCtx := trace.SpanContextFromContext(extractedCtx)
	assert.Equal(t, traceID, extractedSpanCtx.TraceID())
	assert.Equal(t, spanID, extractedSpanCtx.SpanID())
	assert.True(t, extractedSpanCtx.IsSampled())
}

func TestRegisterPropagators_InjectAndExtract(t *testing.T) {
	RegisterPropagators()
	globalProp := otel.GetTextMapPropagator()
	require.NotNil(t, globalProp)

	ctx, traceID, spanID := createTestSpanContext(t)
	header := make(http.Header)
	globalProp.Inject(ctx, propagation.HeaderCarrier(header))
	assert.NotEmpty(t, header.Get("traceparent"))

	extractedCtx := globalProp.Extract(context.Background(), propagation.HeaderCarrier(header))
	assertSpanContextMatches(t, extractedCtx, traceID, spanID)
}
