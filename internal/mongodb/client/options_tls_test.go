package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

const (
	testTLSUser  = "tlsuser"
	testTLSPass  = "tlspass"
	testTLSHost  = "cluster.example.com"
	testUUIDUser = "uuiduser"
	testUUIDPass = "uuidpass"
	testUUIDHost = "uuidhost"
)

// createTLSTestConfig creates a populated test configuration.
func createTLSTestConfig(protocol string) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = testTLSHost
	cfg.Port = testOptionsPort
	cfg.Username = testTLSUser
	cfg.Password = testTLSPass
	cfg.Protocol = protocol
	cfg.UUIDRepresentation = config.UUIDRepresentationStandard
	return cfg
}

// buildUUIDConfig creates a Config with specified UUID representation.
func buildUUIDConfig(rep string) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = testUUIDHost
	cfg.Port = testOptionsPort
	cfg.Username = testUUIDUser
	cfg.Password = testUUIDPass
	cfg.Protocol = config.ProtocolMongoDB
	cfg.UUIDRepresentation = rep
	return cfg
}

// TestBuildClientOptionsWithTLS_Enabled tests options with TLS enabled override.
func TestBuildClientOptionsWithTLS_Enabled(t *testing.T) {
	cfg := createTLSTestConfig(config.ProtocolMongoDB)
	clientOpts := BuildClientOptionsWithTLS(cfg, true)
	require.NotNil(t, clientOpts)

	expectedURI := "mongodb://tlsuser:tlspass@cluster.example.com:27018/?uuidRepresentation=standard&tls=true"
	assert.Equal(t, expectedURI, clientOpts.GetURI())
}

// TestBuildClientOptionsWithTLS_Disabled tests options with TLS disabled override.
func TestBuildClientOptionsWithTLS_Disabled(t *testing.T) {
	cfg := createTLSTestConfig(config.ProtocolMongoDB)
	clientOpts := BuildClientOptionsWithTLS(cfg, false)
	require.NotNil(t, clientOpts)

	expectedURI := "mongodb://tlsuser:tlspass@cluster.example.com:27018/?uuidRepresentation=standard&tls=false"
	assert.Equal(t, expectedURI, clientOpts.GetURI())
}

// TestBuildClientOptions_MongoDBSrv tests options with SRV protocol scheme.
func TestBuildClientOptions_MongoDBSrv(t *testing.T) {
	cfg := createTLSTestConfig(config.ProtocolMongoDBSrv)
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
}

// TestBuildClientOptionsWithTLS_NilConfig tests options built from nil config with TLS flag.
func TestBuildClientOptionsWithTLS_NilConfig(t *testing.T) {
	clientOpts := BuildClientOptionsWithTLS(nil, true)
	require.NotNil(t, clientOpts)
	assert.Equal(t, "mongodb://:27017/?uuidRepresentation=unspecified&tls=true", clientOpts.GetURI())
}

// TestBuildClientOptions_UUIDRepresentations tests options with all UUID representations.
func TestBuildClientOptions_UUIDRepresentations(t *testing.T) {
	testCases := []string{
		config.UUIDRepresentationUnspecified,
		config.UUIDRepresentationStandard,
		config.UUIDRepresentationCSharpLegacy,
		config.UUIDRepresentationJavaLegacy,
		config.UUIDRepresentationPythonLegacy,
	}

	for _, rep := range testCases {
		t.Run(rep, func(t *testing.T) {
			cfg := buildUUIDConfig(rep)
			clientOpts := BuildClientOptions(cfg)
			require.NotNil(t, clientOpts)

			expected := fmt.Sprintf(
				"mongodb://%s:%s@%s:%d/?uuidRepresentation=%s&tls=false",
				testUUIDUser, testUUIDPass, testUUIDHost, testOptionsPort, rep,
			)
			assert.Equal(t, expected, clientOpts.GetURI())
		})
	}
}
