package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testUUIDUser = "uuiduser"
	testUUIDPass = "uuidpass"
	testUUIDHost = "uuidhost"
)

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
