package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_NilConfig(t *testing.T) {
	var cfg *Config
	err := cfg.Validate()
	require.ErrorIs(t, err, ErrNilConfig)
}

func TestValidate_InvalidFields(t *testing.T) {
	t.Run("invalid format", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Format = "yaml"
		require.ErrorIs(t, cfg.Validate(), ErrInvalidFormat)
	})

	t.Run("invalid level", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Level = "verbose"
		require.ErrorIs(t, cfg.Validate(), ErrInvalidLevel)
	})

	t.Run("invalid output", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Output = "unsupported_sink"
		require.ErrorIs(t, cfg.Validate(), ErrInvalidOutput)
	})
}

func TestValidate_FileOutput(t *testing.T) {
	t.Run("invalid file output path", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Output = "/nonexistent_dir_99999/app.log"
		require.ErrorIs(t, cfg.Validate(), ErrInvalidOutput)
	})

	t.Run("valid file output path", func(t *testing.T) {
		logPath := filepath.Join(t.TempDir(), "test.log")
		cfg := DefaultConfig()
		cfg.Output = logPath
		require.NoError(t, cfg.Validate())
	})
}

func TestParsedHelpers_NilConfig(t *testing.T) {
	var cfg *Config

	lvl, err := cfg.ParsedLevel()
	assert.ErrorIs(t, err, ErrNilConfig)
	assert.Empty(t, lvl)

	slvl, err := cfg.ParsedSlogLevel()
	assert.ErrorIs(t, err, ErrNilConfig)
	assert.Equal(t, SlogLevelInfo, slvl)

	fmtEnum, err := cfg.ParsedFormat()
	assert.ErrorIs(t, err, ErrNilConfig)
	assert.Empty(t, fmtEnum)
}
