package config

import (
	"errors"
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
