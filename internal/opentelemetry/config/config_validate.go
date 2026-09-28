package config

import (
	"fmt"
	"math"
	"strings"
)

// isSampleRateBounded checks if the rate is within [MinSampleRate, MaxSampleRate].
func isSampleRateBounded(rate float64) bool {
	if rate < MinSampleRate {
		return false
	}
	return rate <= MaxSampleRate
}

// isSampleRateValid verifies that rate is not NaN, Inf, or out of bounds.
func isSampleRateValid(rate float64) bool {
	if math.IsNaN(rate) {
		return false
	}
	if math.IsInf(rate, 0) {
		return false
	}
	return isSampleRateBounded(rate)
}

// validateSampleRate checks that SampleRate is within [0.0, 1.0] and not NaN or Inf.
func (c *Config) validateSampleRate() error {
	if !isSampleRateValid(c.SampleRate) {
		return fmt.Errorf("%w: got %v", ErrInvalidSampleRate, c.SampleRate)
	}
	return nil
}

// validateServiceName checks that ServiceName is non-empty.
func (c *Config) validateServiceName() error {
	if strings.TrimSpace(c.ServiceName) == "" {
		return ErrEmptyServiceName
	}
	return nil
}

// isValidProtocol checks whether protocol string is supported.
func isValidProtocol(p string) bool {
	return p == ProtocolGRPC || isHTTPProtocol(p)
}

// validateProtocol checks that Protocol is one of the supported OTLP exporter protocols.
func (c *Config) validateProtocol() error {
	p := strings.ToLower(strings.TrimSpace(c.Protocol))
	if !isValidProtocol(p) {
		return fmt.Errorf("%w: %q", ErrInvalidProtocol, c.Protocol)
	}
	return nil
}

// validateActiveConfig validates service name, protocol, and endpoint when enabled.
func (c *Config) validateActiveConfig() error {
	if err := c.validateServiceName(); err != nil {
		return err
	}
	if err := c.validateProtocol(); err != nil {
		return err
	}
	return c.validateEndpoint()
}

// validateEnabledConfig checks active configuration if telemetry is enabled.
func (c *Config) validateEnabledConfig() error {
	if !c.Enabled {
		return nil
	}
	return c.validateActiveConfig()
}

// Validate validates the configuration values.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}
	if err := c.validateSampleRate(); err != nil {
		return err
	}
	return c.validateEnabledConfig()
}
