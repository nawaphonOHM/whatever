package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"github.com/nawaphonOHM/whatever/v3/internal/opentelemetry/config"
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

func TestBuildGRPCOptions_Variants(t *testing.T) {
	cfgWithScheme := &config.Config{
		Endpoint: "dns:///localhost:4317",
		Insecure: true,
	}
	opts1 := buildGRPCOptions(cfgWithScheme)
	assert.NotEmpty(t, opts1)

	cfgWithoutSchemeSecure := &config.Config{
		Endpoint: "localhost:4317",
		Insecure: false,
	}
	opts2 := buildGRPCOptions(cfgWithoutSchemeSecure)
	assert.NotEmpty(t, opts2)
}

func TestBuildHTTPOptions_Variants(t *testing.T) {
	cfgWithSchemeSecure := &config.Config{
		Endpoint: "https://localhost:4318/v1/traces",
		Protocol: config.ProtocolHTTPJSON,
		Insecure: false,
	}
	opts1 := buildHTTPOptions(cfgWithSchemeSecure)
	assert.NotEmpty(t, opts1)

	cfgWithoutSchemeInsecure := &config.Config{
		Endpoint: "localhost:4318",
		Protocol: config.ProtocolHTTPProtobuf,
		Insecure: true,
	}
	opts2 := buildHTTPOptions(cfgWithoutSchemeInsecure)
	assert.NotEmpty(t, opts2)
}
