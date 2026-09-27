package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testTLSUser = "tlsuser"
	testTLSPass = "tlspass"
	testTLSHost = "cluster.example.com"
)

// createTLSTestConfig creates a populated test configuration.
func createTLSTestConfig(protocol string) *config.Config {
	return &config.Config{
		Host:               testTLSHost,
		Port:               testOptionsPort,
		Username:           testTLSUser,
		Password:           testTLSPass,
		Protocol:           protocol,
		UUIDRepresentation: config.UUIDRepresentationStandard,
	}
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
	assert.Empty(t, clientOpts.GetURI())
}
