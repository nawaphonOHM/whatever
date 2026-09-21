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

// ErrNilConfig is returned when operations receive a nil configuration.
var ErrNilConfig = errors.New("mongodb config cannot be nil")

// exitFunc is a package-level hook for os.Exit, allowing tests to intercept process termination.
var exitFunc = os.Exit

// Default configuration constants.
const (
	DefaultConnectTimeout  = 10 * time.Second
	DefaultServerSelection = 5 * time.Second
	DefaultSocketTimeout   = 10 * time.Second
	DefaultMaxPoolSize     = 100
	DefaultMinPoolSize     = 5
	DefaultMaxConnIdleTime = 10 * time.Minute
)

// Config defines configuration options for connecting to MongoDB.
// All environment variable keys use the OHM9996_MONGODB_ prefix.
type Config struct {
	BaseFields    `envPrefix:"OHM9996_MONGODB_"`
	TimeoutFields `envPrefix:"OHM9996_MONGODB_"`
	PoolFields    `envPrefix:"OHM9996_MONGODB_"`
}

// DefaultConfig returns MongoDB configuration with recommended production
// defaults.
func DefaultConfig() *Config {
	return &Config{
		URI:                    "",
		Database:               "",
		Username:               "",
		Password:               "",
		AppName:                "",
		ConnectTimeout:         DefaultConnectTimeout,
		ServerSelectionTimeout: DefaultServerSelection,
		SocketTimeout:          DefaultSocketTimeout,
		MaxConnIdleTime:        DefaultMaxConnIdleTime,
		MaxPoolSize:            DefaultMaxPoolSize,
		MinPoolSize:            DefaultMinPoolSize,
	}
}

// checkMissingRequiredKeys checks whether any required environment keys are unset.
func checkMissingRequiredKeys(cfg *Config) []string {
	checks := []struct {
		val string
		key string
	}{
		{val: cfg.URI, key: "OHM9996_MONGODB_URI"},
		{val: cfg.Database, key: "OHM9996_MONGODB_DATABASE"},
		{val: cfg.Username, key: "OHM9996_MONGODB_USERNAME"},
		{val: cfg.Password, key: "OHM9996_MONGODB_PASSWORD"},
	}

	var missing []string
	for _, c := range checks {
		if c.val == "" {
			missing = append(missing, c.key)
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
