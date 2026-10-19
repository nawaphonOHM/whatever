package provider

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/v3/internal/opentelemetry/config"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
)

const urlSchemeSeparator = "://"

// appendGRPCEndpointOption adds the endpoint configuration to gRPC options.
func appendGRPCEndpointOption(opts []otlptracegrpc.Option, endpoint string) []otlptracegrpc.Option {
	if strings.Contains(endpoint, urlSchemeSeparator) {
		return append(opts, otlptracegrpc.WithEndpointURL(endpoint))
	}
	return append(opts, otlptracegrpc.WithEndpoint(endpoint))
}

// buildGRPCOptions creates options for the gRPC OTLP exporter.
func buildGRPCOptions(cfg *config.Config) []otlptracegrpc.Option {
	opts := appendGRPCEndpointOption(nil, cfg.Endpoint)
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	return opts
}

// newGRPCExporter creates a new OTLP trace exporter over gRPC.
func newGRPCExporter(ctx context.Context, cfg *config.Config) (*otlptrace.Exporter, error) {
	opts := buildGRPCOptions(cfg)
	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc otlp exporter: %w", err)
	}
	return exp, nil
}

// appendHTTPEndpointOption sets endpoint URL or host:port on the option slice.
func appendHTTPEndpointOption(opts []otlptracehttp.Option, endpoint string) []otlptracehttp.Option {
	if strings.Contains(endpoint, urlSchemeSeparator) {
		return append(opts, otlptracehttp.WithEndpointURL(endpoint))
	}
	return append(opts, otlptracehttp.WithEndpoint(endpoint))
}

// resolveHTTPEncoding returns the OTLP HTTP encoding based on protocol configuration.
func resolveHTTPEncoding(protocol string) otlptracehttp.Encoding {
	if strings.EqualFold(strings.TrimSpace(protocol), config.ProtocolHTTPJSON) {
		return otlptracehttp.EncodingJSON
	}
	return otlptracehttp.EncodingProtobuf
}

// buildHTTPOptions creates options for the HTTP OTLP exporter.
func buildHTTPOptions(cfg *config.Config) []otlptracehttp.Option {
	opts := appendHTTPEndpointOption(nil, cfg.Endpoint)
	opts = append(opts, otlptracehttp.WithEncoding(resolveHTTPEncoding(cfg.Protocol)))
	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	return opts
}

// newHTTPExporter creates a new OTLP trace exporter over HTTP.
func newHTTPExporter(ctx context.Context, cfg *config.Config) (*otlptrace.Exporter, error) {
	opts := buildHTTPOptions(cfg)
	exp, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http otlp exporter: %w", err)
	}
	return exp, nil
}

// newExporter constructs an OTLP SpanExporter according to the configured protocol.
func newExporter(ctx context.Context, cfg *config.Config) (sdktrace.SpanExporter, error) {
	if cfg.IsGRPC() {
		logging.InfoContext(ctx, "OpenTelemetry exporter choice",
			"protocol", config.ProtocolGRPC,
			"transport", "grpc",
			"endpoint", cfg.Endpoint,
		)
		return newGRPCExporter(ctx, cfg)
	}
	if cfg.IsHTTP() {
		logging.InfoContext(ctx, "OpenTelemetry exporter choice",
			"protocol", cfg.Protocol,
			"transport", "http",
			"endpoint", cfg.Endpoint,
		)
		return newHTTPExporter(ctx, cfg)
	}
	return nil, fmt.Errorf("%w: %s", config.ErrInvalidProtocol, cfg.Protocol)
}
