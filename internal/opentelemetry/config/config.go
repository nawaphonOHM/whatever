// Package config provides configuration parsing, validation, and defaults
// for OpenTelemetry distributed tracing and telemetry export.
package config

import (
	"fmt"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
	"github.com/nawaphonOHM/whatever/pkg/logging"
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

// logResolution logs the environment file resolution mode.
func logResolution(filenames []string) {
	if len(filenames) == 0 {
		logging.Info("using default OpenTelemetry environment configuration resolution")
		return
	}
	logging.Info("loading OpenTelemetry configuration from files", "files", filenames)
}

// logChoices logs the selected configuration options.
func logChoices(cfg *Config) {
	logging.Info("OpenTelemetry configuration choices",
		"enabled", cfg.Enabled,
		"service_name", cfg.ServiceName,
		"endpoint", cfg.Endpoint,
		"protocol", cfg.Protocol,
		"insecure", cfg.Insecure,
		"sample_rate", cfg.SampleRate,
		"skip_paths", cfg.SkipPaths,
	)
}

// Load loads OpenTelemetry configuration from environment variables,
// applies defaults, and validates the resulting configuration.
// If filenames are provided, environment variables are loaded from those files.
// Otherwise, it checks .env and configs/.env by default.
func Load(filenames ...string) (*Config, error) {
	logging.Info("entering OpenTelemetry configuration initialization state")
	logResolution(filenames)

	cfg, err := intcfg.Load[Config](filenames...)
	if err != nil {
		return nil, fmt.Errorf("failed to load opentelemetry config: %w", err)
	}

	if err := verifyLoadedConfig(cfg); err != nil {
		return nil, err
	}
	logChoices(cfg)
	logging.Info("OpenTelemetry configuration initialization complete")
	return cfg, nil
}
