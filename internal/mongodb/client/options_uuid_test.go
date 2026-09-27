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
	return &config.Config{
		Host:               testUUIDHost,
		Port:               testOptionsPort,
		Username:           testUUIDUser,
		Password:           testUUIDPass,
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: rep,
	}
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
