package provider

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"

	"github.com/nawaphonOHM/whatever/v2/internal/opentelemetry/config"
)

const testHalfSampleRate = 0.5

func verifyActiveProvider(t *testing.T, cfg *config.Config) {
	ctx := context.Background()
	shutdown, err := InitTracerProvider(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	tracer := otel.GetTracerProvider().Tracer("test-tracer")
	assert.NotNil(t, tracer)

	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	assert.NoError(t, shutdown(shutdownCtx))
}

func TestInitTracerProvider_ActiveGRPC(t *testing.T) {
	cfg := &config.Config{
		ServiceName: "grpc-service",
		Endpoint:    "127.0.0.1:4317",
		Protocol:    config.ProtocolGRPC,
		SampleRate:  1.0,
		Enabled:     true,
		Insecure:    true,
	}
	verifyActiveProvider(t, cfg)
}

func TestInitTracerProvider_ActiveHTTPProtobuf(t *testing.T) {
	cfg := &config.Config{
		ServiceName: "http-proto-service",
		Endpoint:    "http://127.0.0.1:4318/v1/traces",
		Protocol:    config.ProtocolHTTPProtobuf,
		SampleRate:  testHalfSampleRate,
		Enabled:     true,
		Insecure:    true,
	}
	verifyActiveProvider(t, cfg)
}

func TestInitTracerProvider_ActiveHTTPJSON(t *testing.T) {
	cfg := &config.Config{
		ServiceName: "http-json-service",
		Endpoint:    "127.0.0.1:4318",
		Protocol:    config.ProtocolHTTPJSON,
		SampleRate:  1.0,
		Enabled:     true,
		Insecure:    true,
	}
	verifyActiveProvider(t, cfg)
}
