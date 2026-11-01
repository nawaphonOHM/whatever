package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	overMaxPort = 65536
	customPort  = 8080
)

// TestConfig_Validate_PortNegative verifies rejection of negative port.
func TestConfig_Validate_PortNegative(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Port = -1
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "port must be between")
}

// TestConfig_Validate_PortOverMax verifies rejection of port exceeding 65535.
func TestConfig_Validate_PortOverMax(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Port = overMaxPort
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "port must be between")
}

// TestConfig_Validate_PortValid verifies acceptance of valid custom port and unset port (0).
func TestConfig_Validate_PortValid(t *testing.T) {
	cfg := validBaseConfig()
	cfg.Port = customPort
	assert.NoError(t, cfg.Validate())

	cfg.Port = 0
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Validate_EndpointInvalid verifies rejection of invalid endpoint URLs.
func TestConfig_Validate_EndpointInvalid(t *testing.T) {
	for _, bad := range []string{
		"http://my endpoint.com",
		"http://",
		"http:///s3path",
		"ftp://minio.local",
		"http://localhost:99999",
		" http://localhost:9000",
		"http://localhost:9000 ",
	} {
		cfg := validBaseConfig()
		cfg.Endpoint = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid endpoint")
	}
}

// TestConfig_Validate_EndpointValid verifies valid endpoint URLs and schemeless hosts.
func TestConfig_Validate_EndpointValid(t *testing.T) {
	for _, good := range []string{
		"http://localhost:9000",
		"https://s3.amazonaws.com",
		"localhost:9000",
		"s3.us-east-1.amazonaws.com",
	} {
		cfg := validBaseConfig()
		cfg.Endpoint = good
		assert.NoError(t, cfg.Validate())
	}
}

// TestConfig_Validate_BucketURLInvalid verifies rejection of invalid bucket URLs.
func TestConfig_Validate_BucketURLInvalid(t *testing.T) {
	for _, bad := range []string{
		"https://bucket url.com",
		"https://",
		"ftp://storage.googleapis.com",
		"https://storage.googleapis.com:99999",
		" https://storage.googleapis.com",
	} {
		cfg := validBaseConfig()
		cfg.BucketURL = bad
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid bucket url")
	}
}

// TestConfig_Validate_ValidEndpoint_InvalidBucketURL verifies that an invalid BucketURL is flagged even when Endpoint is valid.
func TestConfig_Validate_ValidEndpoint_InvalidBucketURL(t *testing.T) {
	for _, badBucket := range []string{
		"https://bucket url.com",
		"https://",
		"http:///path",
		"ftp://storage.googleapis.com",
		"https://storage.googleapis.com/%zz",
		" https://storage.googleapis.com",
	} {
		cfg := validBaseConfig()
		cfg.Endpoint = "https://s3.amazonaws.com"
		cfg.BucketURL = badBucket
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid bucket url")
		assert.NotContains(t, err.Error(), "invalid endpoint")
	}
}

// TestConfig_Validate_CredentialRedaction verifies that credentials in malformed URLs are redacted from validation error messages.
func TestConfig_Validate_CredentialRedaction(t *testing.T) {
	credentialCases := []string{
		"https://AKIAKEY:supersecret@minio.local:9000 /x",
		"https://AKIAKEY:supersecret@minio.local/%zz",
		"https://AKIAKEY:supersecret@",
		"https://AKIAKEY:supersecret@/s3path",
	}

	for _, bad := range credentialCases {
		t.Run("endpoint_"+bad, func(t *testing.T) {
			cfg := validBaseConfig()
			cfg.Endpoint = bad
			err := cfg.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid endpoint")
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")
		})

		t.Run("bucketURL_"+bad, func(t *testing.T) {
			cfg := validBaseConfig()
			cfg.Endpoint = "https://s3.amazonaws.com"
			cfg.BucketURL = bad
			err := cfg.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid bucket url")
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")
		})
	}
}

// TestConfig_Validate_BucketURLValid verifies valid bucket URLs.
func TestConfig_Validate_BucketURLValid(t *testing.T) {
	for _, good := range []string{
		"https://mybucket.s3.us-east-1.amazonaws.com",
		"storage.googleapis.com/mybucket",
	} {
		cfg := validBaseConfig()
		cfg.BucketURL = good
		assert.NoError(t, cfg.Validate())
	}
}
