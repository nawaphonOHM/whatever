package config

import (
	"errors"
	"strings"
)

// Supported OpenTelemetry protocol constants.
const (
	ProtocolGRPC         = "grpc"
	ProtocolHTTP         = "http"
	ProtocolHTTPProtobuf = "http/protobuf"
	ProtocolHTTPJSON     = "http/json"
)

// Boundary constants for port and sample rate.
const (
	MinPort       = 1
	MaxPort       = 65535
	MinSampleRate = 0.0
	MaxSampleRate = 1.0
)

// Sentinel errors for OpenTelemetry configuration validation.
var (
	// ErrNilConfig is returned when validating a nil configuration.
	ErrNilConfig = errors.New("opentelemetry config cannot be nil")

	// ErrEmptyServiceName is returned when the service name is empty while telemetry is enabled.
	ErrEmptyServiceName = errors.New("service name cannot be empty")

	// ErrEmptyEndpoint is returned when the endpoint is empty while telemetry is enabled.
	ErrEmptyEndpoint = errors.New("endpoint cannot be empty")

	// ErrInvalidProtocol is returned when the exporter protocol is unsupported.
	ErrInvalidProtocol = errors.New("invalid protocol: must be grpc, http, http/protobuf, or http/json")

	// ErrInvalidSampleRate is returned when the sample rate is outside [0.0, 1.0] or NaN/Inf.
	ErrInvalidSampleRate = errors.New("sample rate must be between 0.0 and 1.0")

	// ErrInvalidEndpoint is returned when the endpoint URL or host:port is malformed.
	ErrInvalidEndpoint = errors.New("invalid endpoint")
)

// Config defines the configuration options for OpenTelemetry distributed tracing.
// All environment variable keys use the OHM9996_OTEL_ prefix.
type Config struct {
	ServiceName string   `env:"OHM9996_OTEL_SERVICE_NAME" envDefault:"whatever-service"`
	Endpoint    string   `env:"OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4317"`
	Protocol    string   `env:"OHM9996_OTEL_EXPORTER_OTLP_PROTOCOL" envDefault:"grpc"`
	SkipPaths   []string `env:"OHM9996_OTEL_SKIP_PATHS" envSeparator:","`
	SampleRate  float64  `env:"OHM9996_OTEL_SAMPLE_RATE" envDefault:"1.0"`
	Enabled     bool     `env:"OHM9996_OTEL_ENABLED" envDefault:"true"`
	Insecure    bool     `env:"OHM9996_OTEL_INSECURE" envDefault:"true"`
}

// IsGRPC returns true if the configured protocol is gRPC.
func (c *Config) IsGRPC() bool {
	if c == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(c.Protocol), ProtocolGRPC)
}

// isHTTPProtocol returns true if protocol string matches any supported HTTP variant.
func isHTTPProtocol(p string) bool {
	return p == ProtocolHTTP || p == ProtocolHTTPProtobuf || p == ProtocolHTTPJSON
}

// IsHTTP returns true if the configured protocol is any HTTP-based protocol.
func (c *Config) IsHTTP() bool {
	if c == nil {
		return false
	}
	return isHTTPProtocol(strings.ToLower(strings.TrimSpace(c.Protocol)))
}

// isNonEmptyMatch returns true if skip is non-empty and equals path.
func isNonEmptyMatch(skip, path string) bool {
	if skip == "" {
		return false
	}
	return skip == path
}

// matchPath searches skipPaths for a match against path.
func matchPath(skipPaths []string, path string) bool {
	for _, skip := range skipPaths {
		if isNonEmptyMatch(skip, path) {
			return true
		}
	}
	return false
}

// ShouldSkip returns true if the given request path matches any path in SkipPaths.
func (c *Config) ShouldSkip(path string) bool {
	if c == nil {
		return false
	}
	return matchPath(c.SkipPaths, path)
}
