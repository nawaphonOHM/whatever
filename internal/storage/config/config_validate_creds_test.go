package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig_Validate_MissingSecretAccessKey verifies error when only AccessKeyID is set.
func TestConfig_Validate_MissingSecretAccessKey(t *testing.T) {
	cfg := validBaseConfig()
	cfg.AccessKeyID = "my-key-id"
	cfg.SecretAccessKey = ""
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "secret access key cannot be empty")
}

// TestConfig_Validate_MissingAccessKeyID verifies error when only SecretAccessKey is set.
func TestConfig_Validate_MissingAccessKeyID(t *testing.T) {
	cfg := validBaseConfig()
	cfg.AccessKeyID = ""
	cfg.SecretAccessKey = "my-secret-key"
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access key id cannot be empty")
}

// TestConfig_Validate_CredPairValid verifies consistent credential pairs.
func TestConfig_Validate_CredPairValid(t *testing.T) {
	cfg := validBaseConfig()
	cfg.AccessKeyID = "my-key-id"
	cfg.SecretAccessKey = "my-secret"
	assert.NoError(t, cfg.Validate())

	cfg.AccessKeyID = ""
	cfg.SecretAccessKey = ""
	assert.NoError(t, cfg.Validate())
}
