package config

import (
	"os"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Supported environment variable keys for detector.
const (
	EnvCloudProvider                  = "OHM9996_CLOUD_PROVIDER"
	EnvOtelEnabled                    = "OHM9996_OTEL_ENABLED"
	EnvLoggingDisableTraceCorrelation = "OHM9996_LOGGING_DISABLE_TRACE_CORRELATION"
)

// IsGCP reports whether the configured cloud provider is Google Cloud Platform (GCP).
func (c *Config) IsGCP() bool {
	if c == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(c.CloudProvider), CloudProviderGCP)
}

// IsOTelDisabled reports whether OpenTelemetry is disabled for this configuration.
func (c *Config) IsOTelDisabled() bool {
	return IsOTelDisabled(c)
}

// IsGCPStructuredLoggingEnabled reports whether GCP structured logging should be active.
func (c *Config) IsGCPStructuredLoggingEnabled() bool {
	return IsGCPStructuredLoggingEnabled(c)
}

func isTraceCorrelationDisabled(cfg *Config) bool {
	if cfg != nil {
		return cfg.DisableTraceCorrelation
	}
	val, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(EnvLoggingDisableTraceCorrelation)))
	return err == nil && val
}

func parseOtelEnabledEnv() (enabled, ok bool) {
	val, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(EnvOtelEnabled)))
	if err != nil {
		return false, false
	}
	return val, true
}

func isTracerActive() bool {
	tp := otel.GetTracerProvider()
	if tp == nil {
		return false
	}
	_, ok := tp.(*sdktrace.TracerProvider)
	return ok
}

// IsOTelDisabled reports whether OpenTelemetry distributed tracing is disabled or unconfigured.
func IsOTelDisabled(cfg *Config) bool {
	if isTraceCorrelationDisabled(cfg) {
		return true
	}
	if enabled, ok := parseOtelEnabledEnv(); ok {
		return !enabled
	}
	return !isTracerActive()
}

// IsGCPStructuredLoggingEnabled returns true if GCP structured logging should be enabled.
func IsGCPStructuredLoggingEnabled(cfg *Config) bool {
	if cfg == nil {
		return false
	}
	return cfg.IsGCP() && IsOTelDisabled(cfg)
}
