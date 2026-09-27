package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testCustomPortVal = 27018
	testCustomUser    = "testuser"
	testCustomPass    = "testpass"
	testCustomHost    = "custom-host"
)

// buildCustomTestConfig constructs a Config populated with custom parameters.
func buildCustomTestConfig() *config.Config {
	return &config.Config{
		Host:               testCustomHost,
		Port:               testCustomPortVal,
		Username:           testCustomUser,
		Password:           testCustomPass,
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: config.UUIDRepresentationStandard,
	}
}

// TestBuildClientOptions_CustomConfig tests options with explicit values.
func TestBuildClientOptions_CustomConfig(t *testing.T) {
	cfg := buildCustomTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	expectedURI := "mongodb://testuser:testpass@custom-host:27018/?uuidRepresentation=standard&tls=false"
	assert.Equal(t, expectedURI, clientOpts.GetURI())
}

// TestBuildClientOptions_SpecialCharactersAuth tests options with special chars.
func TestBuildClientOptions_SpecialCharactersAuth(t *testing.T) {
	cfg := &config.Config{
		Host:               testCustomHost,
		Port:               testCustomPortVal,
		Username:           "user@name:special",
		Password:           "p@ss/word#123",
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: config.UUIDRepresentationStandard,
	}
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	require.NotNil(t, clientOpts.Auth)
	assert.Equal(t, "user@name:special", clientOpts.Auth.Username)
	assert.Equal(t, "p@ss/word#123", clientOpts.Auth.Password)
}
