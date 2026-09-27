package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig_Validate_UUIDRepresentation tests valid and invalid UUID representation values.
func TestConfig_Validate_UUIDRepresentation(t *testing.T) {
	validReps := []string{
		UUIDRepresentationUnspecified,
		UUIDRepresentationStandard,
		UUIDRepresentationCSharpLegacy,
		UUIDRepresentationJavaLegacy,
		UUIDRepresentationPythonLegacy,
	}

	for _, rep := range validReps {
		cfg := validTestConfig()
		cfg.UUIDRepresentation = rep
		t.Run("valid uuid representation "+rep, func(t *testing.T) {
			assert.NoError(t, cfg.Validate())
		})
	}

	invalidCfg := validTestConfig()
	invalidCfg.UUIDRepresentation = "invalidRep"
	t.Run("invalid uuid representation", func(t *testing.T) {
		err := invalidCfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid mongodb uuid representation: invalidRep")
	})
}
