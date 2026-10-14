package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	overMaxPresignDays    = 8
	testConnectTimeoutSec = 5
	testRequestTimeoutSec = 30
	testPingTimeoutSec    = 5
)

// TestConfig_Validate_ConnectTimeout verifies validation of connect timeout.
func TestConfig_Validate_ConnectTimeout(t *testing.T) {
	for _, bad := range []time.Duration{-1 * time.Second, 0} {
		cfg := validBaseConfig()
		cfg.ConnectTimeout = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "connect timeout must be positive")
	}

	cfg := validBaseConfig()
	cfg.ConnectTimeout = testConnectTimeoutSec * time.Second
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Validate_RequestTimeout verifies validation of request timeout.
func TestConfig_Validate_RequestTimeout(t *testing.T) {
	for _, bad := range []time.Duration{-1 * time.Second, 0} {
		cfg := validBaseConfig()
		cfg.RequestTimeout = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "request timeout must be positive")
	}

	cfg := validBaseConfig()
	cfg.RequestTimeout = testRequestTimeoutSec * time.Second
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Validate_PingTimeout verifies validation of ping timeout.
func TestConfig_Validate_PingTimeout(t *testing.T) {
	for _, bad := range []time.Duration{-1 * time.Second, 0} {
		cfg := validBaseConfig()
		cfg.PingTimeout = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ping timeout must be positive")
	}

	cfg := validBaseConfig()
	cfg.PingTimeout = testPingTimeoutSec * time.Second
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Validate_PresignExpiry_NonPositive verifies rejection of non-positive presign expiry.
func TestConfig_Validate_PresignExpiry_NonPositive(t *testing.T) {
	for _, bad := range []time.Duration{-1 * time.Second, 0} {
		cfg := validBaseConfig()
		cfg.PresignExpiry = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "presign expiry must be positive")
	}

	cfg := validBaseConfig()
	cfg.PresignExpiry = MaxPresignExpiry
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Validate_PresignExpiry_OverMax verifies rejection of presign expiry exceeding 7 days.
func TestConfig_Validate_PresignExpiry_OverMax(t *testing.T) {
	cfg := validBaseConfig()
	cfg.PresignExpiry = overMaxPresignDays * hoursPerDay * time.Hour
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot exceed")
}
