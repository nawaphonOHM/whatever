package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDefaultConnectTimeout = 10 * time.Second
	testDefaultRequestTimeout = 30 * time.Second
	testDefaultPingTimeout    = 5 * time.Second
	testDefaultPresignExpiry  = 15 * time.Minute
)

func assertDefaultTimeouts(t *testing.T, cfg *Config) {
	assert.Equal(t, testDefaultConnectTimeout, cfg.ConnectTimeout)
	assert.Equal(t, testDefaultRequestTimeout, cfg.RequestTimeout)
	assert.Equal(t, testDefaultPingTimeout, cfg.PingTimeout)
	assert.Equal(t, testDefaultPresignExpiry, cfg.PresignExpiry)
}

func assertDefaultEndpoints(t *testing.T, cfg *Config) {
	assert.Equal(t, 0, cfg.Port)
	assert.Empty(t, cfg.Endpoint)
	assert.Empty(t, cfg.BucketURL)
	assert.False(t, cfg.UsePathStyle)
	assert.False(t, cfg.Insecure)
	assert.False(t, cfg.DisableSSL)
}

// TestDefaultConfig verifies default values on new config instance.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, "us-east-1", cfg.Region)
	assert.Empty(t, cfg.Bucket)
	assertDefaultEndpoints(t, cfg)
	assertDefaultTimeouts(t, cfg)
}

// TestConfig_SetDefaults verifies SetDefaults populates recommended values.
func TestConfig_SetDefaults(t *testing.T) {
	var cfg Config
	cfg.SetDefaults()
	assert.Equal(t, "us-east-1", cfg.Region)
	assertDefaultEndpoints(t, &cfg)
	assertDefaultTimeouts(t, &cfg)
}
