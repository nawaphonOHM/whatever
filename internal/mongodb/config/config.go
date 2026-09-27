// Package config provides configuration parsing, validation, and defaults
// for MongoDB connections.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
)

const (
	defaultPort                   = 27017
	defaultConnectTimeoutSec      = 10
	defaultServerSelectionTimeout = 5
	defaultSocketTimeoutSec       = 10
	defaultMaxConnIdleMinutes     = 10
	defaultMaxPoolSize            = 100
	defaultMinPoolSize            = 5
)

// ErrNilConfig is returned when operations receive a nil configuration.
var ErrNilConfig = errors.New("mongodb config cannot be nil")

// exitFunc is a package-level hook for os.Exit, allowing tests to intercept process termination.
var exitFunc = os.Exit

// SetDefaults populates the configuration with initial default values.
func (c *Config) SetDefaults() {
	*c = *DefaultConfig()
}

// DefaultConfig returns MongoDB configuration with recommended defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:                   defaultPort,
		Protocol:               ProtocolMongoDB,
		UUIDRepresentation:     UUIDRepresentationUnspecified,
		ConnectTimeout:         defaultConnectTimeoutSec * time.Second,
		ServerSelectionTimeout: defaultServerSelectionTimeout * time.Second,
		SocketTimeout:          defaultSocketTimeoutSec * time.Second,
		MaxConnIdleTime:        defaultMaxConnIdleMinutes * time.Minute,
		MaxPoolSize:            defaultMaxPoolSize,
		MinPoolSize:            defaultMinPoolSize,
	}
}

type requiredKeyCheck struct {
	key     string
	missing bool
}

// buildRequiredKeyChecks returns presence checks for mandatory configuration keys.
func (c *Config) buildRequiredKeyChecks() []requiredKeyCheck {
	return []requiredKeyCheck{
		{key: "OHM9996_MONGODB_HOST", missing: c.Host == ""},
	}
}

// checkMissingRequiredKeys checks whether any required environment keys are unset.
func checkMissingRequiredKeys(cfg *Config) []string {
	var missing []string
	for _, check := range cfg.buildRequiredKeyChecks() {
		if check.missing {
			missing = append(missing, check.key)
		}
	}
	return missing
}

const missingKeysLogFormat = "missing required mongodb configuration keys: %v; exiting peacefully\n"

// logMissingKeys outputs diagnostic message for missing environment keys.
func logMissingKeys(missing []string) {
	if _, err := fmt.Fprintf(os.Stderr, missingKeysLogFormat, missing); err != nil {
		return
	}
}

// handleMissingKeys inspects missing keys and initiates peaceful termination if any are missing.
func handleMissingKeys(missing []string) error {
	if len(missing) == 0 {
		return nil
	}
	logMissingKeys(missing)
	exitFunc(0)
	return fmt.Errorf("missing required mongodb configuration keys: %v", missing)
}

// verifyLoadedConfig ensures required environment keys are present and settings are valid.
func verifyLoadedConfig(cfg *Config) error {
	if err := handleMissingKeys(checkMissingRequiredKeys(cfg)); err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid mongodb config: %w", err)
	}
	return nil
}

// LoadConfig loads MongoDB configuration from environment variables with
// OHM9996_MONGODB_* prefix and validates the resulting settings.
// If any required connection or credential settings are missing, it logs
// a diagnostic message and terminates peacefully with exit code 0.
func LoadConfig(filenames ...string) (*Config, error) {
	cfg, err := intcfg.Load[Config](filenames...)
	if err != nil {
		return nil, fmt.Errorf("failed to load mongodb config: %w", err)
	}

	if err := verifyLoadedConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
