package provider

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
)

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
