// Package config provides configuration parsing, validation, and defaults
// for OpenTelemetry distributed tracing and telemetry export.
package config

import (
	"fmt"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
)

// Default configuration constants.
const (
	DefaultEnabled     = true
	DefaultServiceName = "whatever-service"
	DefaultEndpoint    = "localhost:4317"
	DefaultProtocol    = ProtocolGRPC
	DefaultInsecure    = true
	DefaultSampleRate  = 1.0
)

// SetDefaults populates the configuration with initial default values.
func (c *Config) SetDefaults() {
	*c = *DefaultConfig()
}

// DefaultConfig returns OpenTelemetry configuration with recommended defaults.
func DefaultConfig() *Config {
	return &Config{
		SkipPaths:   nil,
		ServiceName: DefaultServiceName,
		Endpoint:    DefaultEndpoint,
		Protocol:    DefaultProtocol,
		SampleRate:  DefaultSampleRate,
		Enabled:     DefaultEnabled,
		Insecure:    DefaultInsecure,
	}
}

// verifyLoadedConfig runs validation on the loaded configuration.
func verifyLoadedConfig(cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid opentelemetry config: %w", err)
	}
	return nil
}

// Load loads OpenTelemetry configuration from environment variables,
// applies defaults, and validates the resulting configuration.
// If filenames are provided, environment variables are loaded from those files.
// Otherwise, it checks .env and configs/.env by default.
func Load(filenames ...string) (*Config, error) {
	cfg, err := intcfg.Load[Config](filenames...)
	if err != nil {
		return nil, fmt.Errorf("failed to load opentelemetry config: %w", err)
	}

	if err := verifyLoadedConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
