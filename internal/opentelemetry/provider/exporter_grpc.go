package provider

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"github.com/nawaphonOHM/whatever/v2/internal/opentelemetry/config"
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
