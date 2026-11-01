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
		"https:/admin:SECRETPW@minio.internal:9000",
		"//admin:SECRETPW@minio.internal:9000",
		"https://admin:SECRETPW@minio.internal:SECRETPW",
		"minio.local/a/b://user:SECRETPW@host",
	} {
		cfg := validBaseConfig()
		cfg.Endpoint = bad
		require.Error(t, cfg.Validate())
	}
}

// TestConfig_Validate_EndpointValid verifies valid endpoint URLs and schemeless hosts.
func TestConfig_Validate_EndpointValid(t *testing.T) {
	for _, good := range []string{
		"http://localhost:9000",
		"https://s3.amazonaws.com",
		"localhost:9000",
		"s3.us-east-1.amazonaws.com",
		"http://[::1]:9000",
		"[::1]:9000",
		"https://[::1]/x",
		"http://[::ffff:127.0.0.1]:9000",
		"https://minio.local/tenant@team",
		"https://proxy.example.com/s3/a@b/",
		"https://minio.local/a%40b",
	} {
		cfg := validBaseConfig()
		cfg.Endpoint = good
		assert.NoError(t, cfg.Validate(), good)
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
		require.Error(t, cfg.Validate())
	}
}

// TestConfig_Validate_BucketURLValid verifies valid bucket URLs.
func TestConfig_Validate_BucketURLValid(t *testing.T) {
	for _, good := range []string{
		"https://mybucket.s3.us-east-1.amazonaws.com",
		"storage.googleapis.com/mybucket",
		"https://mybucket.s3.amazonaws.com/key@v1",
		"https://minio.local/a%40b",
	} {
		cfg := validBaseConfig()
		cfg.BucketURL = good
		assert.NoError(t, cfg.Validate(), good)
	}
}
