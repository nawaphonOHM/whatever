package logging

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validConfigCase struct {
	name string
	cfg  Config
}

func validConfigCases() []validConfigCase {
	return []validConfigCase{
		{
			name: "valid json config with debug level",
			cfg: Config{
				Output: &bytes.Buffer{},
				Level:  LevelDebug,
				Format: FormatJSON,
			},
		},
		{
			name: "valid text config with trace level",
			cfg: Config{
				Output: &bytes.Buffer{},
				Level:  LevelTrace,
				Format: FormatText,
			},
		},
		{
			name: "valid config with fatal level",
			cfg: Config{
				Output: &bytes.Buffer{},
				Level:  LevelFatal,
				Format: FormatJSON,
			},
		},
		{
			name: "empty format and level defaults",
			cfg: Config{
				Output: &bytes.Buffer{},
			},
		},
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, LevelInfo, cfg.Level)
	assert.Equal(t, FormatJSON, cfg.Format)
	assert.Equal(t, os.Stdout, cfg.Output)
	assert.False(t, cfg.AddSource)
	assert.False(t, cfg.DisableTraceCorrelation)
	assert.NotNil(t, cfg.ExitFunc)
	assert.NoError(t, cfg.Validate())
}

func TestConfig_Validate_Valid(t *testing.T) {
	for _, tt := range validConfigCases() {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.cfg.Validate())
		})
	}
}

func TestConfig_Validate_Invalid(t *testing.T) {
	t.Run("invalid format", func(t *testing.T) {
		cfg := Config{Format: "yaml"}
		err := cfg.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFormat)
	})

	t.Run("invalid level", func(t *testing.T) {
		cfg := Config{Level: "invalid-level"}
		err := cfg.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidLevel)
	})
}
