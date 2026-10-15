package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// TestConfig_Validate_NilConfig verifies nil pointer error matches sentinel errors across packages.
func TestConfig_Validate_NilConfig(t *testing.T) {
	var cfg *config.Config
	err := cfg.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrNilConfig))
	assert.True(t, errors.Is(err, contracts.ErrNilConfig))
}
