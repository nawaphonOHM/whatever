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
		"https://proxy.company.com/team%40s3",
		"https://proxy.company.com/team/s3",
		"https://user:p%40ss@minio.local:9000",
		"https://user:p@ss@minio.local:9000",
		"https://user:p@ss@minio.local:9000/s3",
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
	credentialCases := []struct {
		name      string
		target    string
		forbidden []string
	}{
		{
			name:      "alphanumeric credentials with whitespace",
			target:    "https://AKIAKEY:supersecret@minio.local:9000 /x",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with bad escape",
			target:    "https://AKIAKEY:supersecret@minio.local/%zz",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with missing host",
			target:    "https://AKIAKEY:supersecret@",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with missing host and path",
			target:    "https://AKIAKEY:supersecret@/s3path",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "AWS-style secret with slashes and whitespace",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@minio.local:9000 /x",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and bad escape",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@minio.local/%zz",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and missing host",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and missing host with path",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@/s3path",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and trailing at without host",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG@",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG"},
		},
		{
			name:      "secret containing hash with trailing at",
			target:    "https://AKIAKEY:hash#secret@",
			forbidden: []string{"AKIAKEY", "hash#secret", "hash", "secret"},
		},
		{
			name:      "secret containing hash with whitespace",
			target:    "https://AKIAKEY:hash#secret@minio.local:9000 /x",
			forbidden: []string{"AKIAKEY", "hash#secret", "hash", "secret"},
		},
		{
			name:      "secret containing question mark with trailing at",
			target:    "https://AKIAKEY:quest?secret@",
			forbidden: []string{"AKIAKEY", "quest?secret", "quest", "secret"},
		},
		{
			name:      "secret containing question mark with bad escape",
			target:    "https://AKIAKEY:quest?secret@minio.local/%zz",
			forbidden: []string{"AKIAKEY", "quest?secret", "quest", "secret"},
		},
		{
			name:      "unencoded userinfo separator with path",
			target:    "https://user:8080/supersecret/key@s3.amazonaws.com",
			forbidden: []string{"supersecret", "key@s3.amazonaws.com"},
		},
		{
			name:      "unencoded userinfo separator with query",
			target:    "https://user:8080?x=secret@s3.amazonaws.com",
			forbidden: []string{"secret", "s3.amazonaws.com"},
		},
		{
			name:      "unencoded userinfo separator with fragment",
			target:    "https://user:8080#frag_secret@s3.amazonaws.com",
			forbidden: []string{"frag_secret", "s3.amazonaws.com"},
		},
		{
			name:      "schemeless unencoded userinfo separator with path",
			target:    "user:8080/secret@s3.amazonaws.com",
			forbidden: []string{"secret", "s3.amazonaws.com"},
		},
		{
			name:      "AWS-style secret with slashes before endpoint host",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@s3.amazonaws.com",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
	}

	for _, tc := range credentialCases {
		t.Run("endpoint_"+tc.name, func(t *testing.T) {
			cfg := validBaseConfig()
			cfg.Endpoint = tc.target
			err := cfg.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid endpoint")
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}
		})

		t.Run("bucketURL_"+tc.name, func(t *testing.T) {
			cfg := validBaseConfig()
			cfg.Endpoint = "https://s3.amazonaws.com"
			cfg.BucketURL = tc.target
			err := cfg.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid bucket url")
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}
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
