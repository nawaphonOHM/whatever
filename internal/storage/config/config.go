// Package config provides configuration parsing, validation, and defaults
// for object storage connections.
package config

import (
	"context"
	"fmt"
	"time"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	intcfg "github.com/nawaphonOHM/whatever/v3/internal/rest/config"
)

const (
	defaultConnectTimeoutSec = 10
	defaultRequestTimeoutSec = 30
	defaultPingTimeoutSec    = 5
	defaultPresignExpiryMin  = 15
)

// SetExitFunc overrides the process exit hook used by the central logger worker
// and returns the previous hook.
func SetExitFunc(fn func(int)) func(int) {
	// False positive: borrowing the central worker singleton to configure exit hook is leak-free;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	prev := central.DefaultWorker().ExitFunc()
	central.DefaultWorker().SetExitFunc(fn)
	return prev
}

// SetDefaults populates the configuration with initial default values.
func (c *Config) SetDefaults() {
	*c = *DefaultConfig()
}

// DefaultConfig returns storage configuration with recommended defaults.
func DefaultConfig() *Config {
	var cfg Config
	cfg.Region = "us-east-1"
	cfg.ConnectTimeout = defaultConnectTimeoutSec * time.Second
	cfg.RequestTimeout = defaultRequestTimeoutSec * time.Second
	cfg.PingTimeout = defaultPingTimeoutSec * time.Second
	cfg.PresignExpiry = defaultPresignExpiryMin * time.Minute
	return &cfg
}

type requiredKeyCheck struct {
	key     string
	missing bool
}

// buildRequiredKeyChecks returns presence checks for mandatory configuration keys.
func (c *Config) buildRequiredKeyChecks() []requiredKeyCheck {
	return []requiredKeyCheck{
		{key: "OHM9996_STORAGE_BUCKET", missing: c.Bucket == ""},
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

// handleMissingKeys inspects missing keys and initiates peaceful termination via Central Log if any are missing.
func handleMissingKeys(missing []string) error {
	if len(missing) == 0 {
		return nil
	}
	msg := fmt.Sprintf("missing required storage configuration keys: %v; exiting peacefully", missing)
	central.ExitWithGraceful(context.Background(), msg)
	return fmt.Errorf("missing required storage configuration keys: %v", missing)
}

// verifyLoadedConfig ensures required environment keys are present and settings are valid.
func verifyLoadedConfig(cfg *Config) error {
	if err := handleMissingKeys(checkMissingRequiredKeys(cfg)); err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid storage config: %w", err)
	}
	return nil
}

// LoadConfig loads storage configuration from environment variables with
// OHM9996_STORAGE_* prefix (with fallback to .env/configs/.env via rest/config.Load)
// and validates the resulting settings.
// If any required connection settings are missing, it logs
// a diagnostic message and terminates peacefully with exit code 0.
func LoadConfig(filenames ...string) (*Config, error) {
	cfg, err := intcfg.Load[Config](filenames...)
	if err != nil {
		return nil, fmt.Errorf("failed to load storage config: %w", err)
	}

	if err := verifyLoadedConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
