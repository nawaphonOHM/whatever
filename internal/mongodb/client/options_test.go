package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testOptionsPort = 27018
)

// TestBuildClientOptions_DefaultConfig tests options built from DefaultConfig.
func TestBuildClientOptions_DefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	assert.Empty(t, clientOpts.GetURI())
}

// TestBuildClientOptions_NilConfig tests options built when config is nil.
func TestBuildClientOptions_NilConfig(t *testing.T) {
	clientOpts := BuildClientOptions(nil)
	require.NotNil(t, clientOpts)
	assert.Empty(t, clientOpts.GetURI())
}

// TestBuildClientOptions_PopulatedConfig tests options with populated config.
func TestBuildClientOptions_PopulatedConfig(t *testing.T) {
	cfg := &config.Config{
		Host:               "custom-host",
		Port:               testOptionsPort,
		Username:           "admin",
		Password:           "secretpassword",
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: config.UUIDRepresentationStandard,
	}
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	expectedURI := "mongodb://admin:secretpassword@custom-host:27018/?uuidRepresentation=standard&tls=false"
	assert.Equal(t, expectedURI, clientOpts.GetURI())
}

// TestBuildClientOptions_AuthAndHosts verifies hosts and credentials on options.
func TestBuildClientOptions_AuthAndHosts(t *testing.T) {
	cfg := &config.Config{
		Host:               "custom-host",
		Port:               testOptionsPort,
		Username:           "admin",
		Password:           "secretpassword",
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: config.UUIDRepresentationStandard,
	}
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	assert.Contains(t, clientOpts.Hosts, "custom-host:27018")
	require.NotNil(t, clientOpts.Auth)
	assert.Equal(t, "admin", clientOpts.Auth.Username)
	assert.Equal(t, "secretpassword", clientOpts.Auth.Password)
}
