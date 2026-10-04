package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

const (
	testCustomPortVal = 27018
	testCustomUser    = "testuser"
	testCustomPass    = "testpass"
	testCustomHost    = "custom-host"
)

// buildCustomTestConfig constructs a Config populated with custom parameters.
func buildCustomTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = testCustomHost
	cfg.Port = testCustomPortVal
	cfg.Username = testCustomUser
	cfg.Password = testCustomPass
	cfg.Protocol = config.ProtocolMongoDB
	cfg.UUIDRepresentation = config.UUIDRepresentationStandard
	return cfg
}

// TestBuildClientOptions_CustomConfig tests options with explicit values.
func TestBuildClientOptions_CustomConfig(t *testing.T) {
	cfg := buildCustomTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	expectedURI := "mongodb://testuser:testpass@custom-host:27018/?uuidRepresentation=standard&tls=false"
	assert.Equal(t, expectedURI, clientOpts.GetURI())
}

// buildSpecialAuthTestConfig constructs a Config with special auth characters.
func buildSpecialAuthTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = testCustomHost
	cfg.Port = testCustomPortVal
	cfg.Username = "user@name:special"
	cfg.Password = "p@ss/word#123"
	cfg.Protocol = config.ProtocolMongoDB
	cfg.UUIDRepresentation = config.UUIDRepresentationStandard
	return cfg
}

// TestBuildClientOptions_SpecialCharactersAuth tests options with special chars.
func TestBuildClientOptions_SpecialCharactersAuth(t *testing.T) {
	cfg := buildSpecialAuthTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	require.NotNil(t, clientOpts.Auth)
	assert.Equal(t, "user@name:special", clientOpts.Auth.Username)
	assert.Equal(t, "p@ss/word#123", clientOpts.Auth.Password)
}
