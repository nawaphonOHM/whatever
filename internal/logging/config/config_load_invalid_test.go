package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_InvalidFormat(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvFormat, "xml")

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.ErrorIs(t, err, ErrInvalidFormat)
}

func TestLoad_InvalidLevel(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvLevel, "invalid_level")

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.ErrorIs(t, err, ErrInvalidLevel)
}

func TestLoad_InvalidOutput(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvOutput, "invalid_dest_name")

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.ErrorIs(t, err, ErrInvalidOutput)
}

func TestLoad_InvalidTypeConversion(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvAddSource, "not-a-bool")

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
}
