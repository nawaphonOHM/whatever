package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"github.com/nawaphonOHM/whatever/v2/internal/opentelemetry/config"
)

func TestResolveHTTPEncoding(t *testing.T) {
	assert.Equal(t, otlptracehttp.EncodingJSON, resolveHTTPEncoding(config.ProtocolHTTPJSON))
	assert.Equal(t, otlptracehttp.EncodingProtobuf, resolveHTTPEncoding(config.ProtocolHTTP))
	assert.Equal(t, otlptracehttp.EncodingProtobuf, resolveHTTPEncoding(config.ProtocolHTTPProtobuf))
	assert.Equal(t, otlptracehttp.EncodingProtobuf, resolveHTTPEncoding("other"))
}

func TestNewExporter_InvalidProtocol(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		ServiceName: "test-service",
		Endpoint:    "localhost:4317",
		Protocol:    "unsupported",
		Enabled:     true,
	}

	exp, err := newExporter(ctx, cfg)
	require.Error(t, err)
	assert.Nil(t, exp)
	assert.True(t, errors.Is(err, config.ErrInvalidProtocol))
}

func TestNewExporter_Success(t *testing.T) {
	ctx := context.Background()

	grpcCfg := &config.Config{
		ServiceName: "grpc-svc",
		Endpoint:    "127.0.0.1:4317",
		Protocol:    config.ProtocolGRPC,
		Insecure:    true,
		Enabled:     true,
	}
	grpcExp, err := newExporter(ctx, grpcCfg)
	require.NoError(t, err)
	assert.NotNil(t, grpcExp)

	httpCfg := &config.Config{
		ServiceName: "http-svc",
		Endpoint:    "http://127.0.0.1:4318",
		Protocol:    config.ProtocolHTTP,
		Insecure:    true,
		Enabled:     true,
	}
	httpExp, err := newExporter(ctx, httpCfg)
	require.NoError(t, err)
	assert.NotNil(t, httpExp)
}
