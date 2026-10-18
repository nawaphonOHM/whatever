package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPort              = 9000
	testPortStr           = "9000"
	testEnvConnectTimeout = 15 * time.Second
	testEnvRequestTimeout = 45 * time.Second
	testEnvPingTimeout    = 8 * time.Second
	testEnvPresignExpiry  = 20 * time.Minute
)

func setupFullAuthEnv(t *testing.T) {
	t.Setenv("OHM9996_STORAGE_BUCKET", "test-bucket")
	t.Setenv("OHM9996_STORAGE_REGION", "ap-southeast-1")
	t.Setenv("OHM9996_STORAGE_ACCESS_KEY_ID", "test-key-id")
	t.Setenv("OHM9996_STORAGE_SECRET_ACCESS_KEY", "test-secret-key")
	t.Setenv("OHM9996_STORAGE_SESSION_TOKEN", "test-token")
}

func setupFullNetworkEnv(t *testing.T) {
	t.Setenv("OHM9996_STORAGE_ENDPOINT", "http://localhost:9000")
	t.Setenv("OHM9996_STORAGE_BUCKET_URL", "http://localhost:9000/test-bucket")
	t.Setenv("OHM9996_STORAGE_PORT", testPortStr)
	t.Setenv("OHM9996_STORAGE_USE_PATH_STYLE", "true")
	t.Setenv("OHM9996_STORAGE_INSECURE", "true")
	t.Setenv("OHM9996_STORAGE_DISABLE_SSL", "true")
	t.Setenv("OHM9996_STORAGE_CONNECT_TIMEOUT", "15s")
	t.Setenv("OHM9996_STORAGE_REQUEST_TIMEOUT", "45s")
	t.Setenv("OHM9996_STORAGE_PING_TIMEOUT", "8s")
	t.Setenv("OHM9996_STORAGE_PRESIGN_EXPIRY", "20m")
}

func assertLoadedAuthFields(t *testing.T, cfg *Config) {
	assert.Equal(t, "test-bucket", cfg.Bucket)
	assert.Equal(t, "ap-southeast-1", cfg.Region)
	assert.Equal(t, "test-key-id", cfg.AccessKeyID)
	assert.Equal(t, "test-secret-key", cfg.SecretAccessKey)
	assert.Equal(t, "test-token", cfg.SessionToken)
}

func assertLoadedNetworkFields(t *testing.T, cfg *Config) {
	assert.Equal(t, "http://localhost:9000", cfg.Endpoint)
	assert.Equal(t, "http://localhost:9000/test-bucket", cfg.BucketURL)
	assert.Equal(t, testPort, cfg.Port)
	assert.True(t, cfg.UsePathStyle)
	assert.True(t, cfg.Insecure)
	assert.True(t, cfg.DisableSSL)
	assert.Equal(t, testEnvConnectTimeout, cfg.ConnectTimeout)
	assert.Equal(t, testEnvRequestTimeout, cfg.RequestTimeout)
	assert.Equal(t, testEnvPingTimeout, cfg.PingTimeout)
	assert.Equal(t, testEnvPresignExpiry, cfg.PresignExpiry)
}

// TestConfig_Load_FullEnv verifies loading all environment variables.
func TestConfig_Load_FullEnv(t *testing.T) {
	setupFullAuthEnv(t)
	setupFullNetworkEnv(t)
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assertLoadedAuthFields(t, cfg)
	assertLoadedNetworkFields(t, cfg)
}

// TestConfig_Load_MinimalEnv verifies loading with only required bucket set.
func TestConfig_Load_MinimalEnv(t *testing.T) {
	t.Setenv("OHM9996_STORAGE_BUCKET", "minimal-bucket")
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "minimal-bucket", cfg.Bucket)
	assert.Equal(t, "us-east-1", cfg.Region)
	assertDefaultEndpoints(t, cfg)
	assertDefaultTimeouts(t, cfg)
}
