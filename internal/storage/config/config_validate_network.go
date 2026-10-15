package config

import (
	"errors"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

// checkPortRange checks that the port is within the valid range (1-65535) or 0 (unset).
func (c *Config) checkPortRange() error {
	if c.Port < 0 || c.Port > MaxPort {
		return fmt.Errorf("storage port must be between %d and %d or 0", MinPort, MaxPort)
	}
	return nil
}

// validateURLFormat parses and verifies URL scheme, host, and port via detector.ParseTarget.
func validateURLFormat(rawURL, fieldName string) error {
	if rawURL == "" {
		return nil
	}
	if _, err := detector.ParseTarget(rawURL); err != nil {
		return fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	return nil
}

// validateConnectionTimeouts checks connect and request timeouts are strictly positive.
func (c *Config) validateConnectionTimeouts() error {
	if c.ConnectTimeout <= 0 {
		return errors.New("storage connect timeout must be positive (> 0)")
	}
	if c.RequestTimeout <= 0 {
		return errors.New("storage request timeout must be positive (> 0)")
	}
	return nil
}

// validatePingTimeout checks ping timeout is strictly positive.
func (c *Config) validatePingTimeout() error {
	if c.PingTimeout <= 0 {
		return errors.New("storage ping timeout must be positive (> 0)")
	}
	return nil
}

// validatePresignExpiry checks presign expiry is strictly positive and within maximum limit.
func (c *Config) validatePresignExpiry() error {
	if c.PresignExpiry <= 0 {
		return errors.New("storage presign expiry must be positive (> 0)")
	}
	if c.PresignExpiry > MaxPresignExpiry {
		return fmt.Errorf("storage presign expiry cannot exceed %v (7 days)", MaxPresignExpiry)
	}
	return nil
}

// validateOperationTimeouts checks ping timeout and presign expiry.
func (c *Config) validateOperationTimeouts() error {
	if err := c.validatePingTimeout(); err != nil {
		return err
	}
	return c.validatePresignExpiry()
}

// validateTimeouts ensures all configured timeout durations are positive and within bounds.
func (c *Config) validateTimeouts() error {
	if err := c.validateConnectionTimeouts(); err != nil {
		return err
	}
	return c.validateOperationTimeouts()
}

// validateEndpoints verifies endpoint and bucket URL formats.
func (c *Config) validateEndpoints() error {
	if err := validateURLFormat(c.Endpoint, "endpoint"); err != nil {
		return err
	}
	return validateURLFormat(c.BucketURL, "bucket url")
}

// validateSettings verifies port, endpoints, and timeout settings.
func (c *Config) validateSettings() error {
	if err := c.checkPortRange(); err != nil {
		return err
	}
	if err := c.validateEndpoints(); err != nil {
		return err
	}
	return c.validateTimeouts()
}
