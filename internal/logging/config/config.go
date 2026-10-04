// Package config provides environment variable parsing, validation, and defaults
// for application logging.
package config

import (
	"fmt"
	"log/slog"

	intcfg "github.com/nawaphonOHM/whatever/v2/internal/rest/config"
)

// SetDefaults populates the configuration with initial default values.
func (c *Config) SetDefaults() {
	*c = *DefaultConfig()
}

// DefaultConfig returns logging configuration with recommended defaults.
func DefaultConfig() *Config {
	return &Config{
		Output:                  DefaultOutput,
		Level:                   DefaultLevel,
		Format:                  DefaultFormat,
		AddSource:               DefaultAddSource,
		DisableTraceCorrelation: DefaultDisableTraceCorrelation,
	}
}

// ParsedLevel returns the parsed Level enum from Config.Level.
func (c *Config) ParsedLevel() (Level, error) {
	if c == nil {
		return "", ErrNilConfig
	}
	return ParseLevel(c.Level)
}

// ParsedSlogLevel returns the parsed slog.Level from Config.Level.
func (c *Config) ParsedSlogLevel() (slog.Level, error) {
	if c == nil {
		return SlogLevelInfo, ErrNilConfig
	}
	return ParseSlogLevel(c.Level)
}

// ParsedFormat returns the parsed Format enum from Config.Format.
func (c *Config) ParsedFormat() (Format, error) {
	if c == nil {
		return "", ErrNilConfig
	}
	return ParseFormat(c.Format)
}

// verifyLoadedConfig runs validation on the loaded configuration.
func verifyLoadedConfig(cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid logging config: %w", err)
	}
	return nil
}

// Load loads logging configuration from environment variables,
// applies defaults, and validates the resulting configuration.
// If filenames are provided, environment variables are loaded from those files.
// Otherwise, it checks .env and configs/.env by default.
func Load(filenames ...string) (*Config, error) {
	cfg, err := intcfg.Load[Config](filenames...)
	if err != nil {
		return nil, fmt.Errorf("failed to load logging config: %w", err)
	}

	if err := verifyLoadedConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
