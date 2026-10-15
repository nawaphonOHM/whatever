package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	maxBucketLength     = 222
	overMaxBucketLength = 223
)

func validBaseConfig() *Config {
	cfg := DefaultConfig()
	cfg.Bucket = "valid-bucket-name"
	return cfg
}

// TestConfig_Validate_BucketEmpty verifies empty bucket is rejected.
func TestConfig_Validate_BucketEmpty(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Bucket = ""
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "storage bucket cannot be empty")
}

// TestConfig_Validate_BucketLength verifies bucket length boundary enforcement.
func TestConfig_Validate_BucketLength(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Bucket = "ab"
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "length must be between")

	cfg.Bucket = strings.Repeat("a", maxBucketLength)
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Validate_BucketLength_OverMax verifies rejection when bucket length exceeds maximum.
func TestConfig_Validate_BucketLength_OverMax(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Bucket = strings.Repeat("a", overMaxBucketLength)
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "length must be between")
}

// TestConfig_Validate_BucketChars verifies forbidden characters in bucket name.
func TestConfig_Validate_BucketChars(t *testing.T) {
	for _, bad := range []string{"my bucket", "my/bucket", "my\\bucket", "my\tbucket"} {
		cfg := validBaseConfig()
		cfg.Bucket = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "contains invalid characters")
	}
}

// TestConfig_Validate_BucketValid verifies compliant bucket names pass.
func TestConfig_Validate_BucketValid(t *testing.T) {
	for _, good := range []string{"abc", "my-bucket", "test.bucket.123", "a-1-b"} {
		cfg := validBaseConfig()
		cfg.Bucket = good
		assert.NoError(t, cfg.Validate())
	}
}
