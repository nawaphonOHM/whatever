package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSmallPoolMax = uint64(5)
	testMinPoolLimit = uint64(10)
	testZeroPoolMax  = uint64(0)
)

// configValidateCase defines a test case for Config.Validate.
type configValidateCase struct {
	cfg         *Config
	name        string
	expectedErr string
}

// buildValidateCases creates table-driven validation cases.
func buildValidateCases() []configValidateCase {
	return []configValidateCase{
		{
			cfg:         nil,
			name:        "nil config",
			expectedErr: "mongodb config cannot be nil",
		},
		{
			cfg:         &Config{URI: ""},
			name:        "empty URI",
			expectedErr: "mongodb uri cannot be empty",
		},
		{
			cfg: &Config{
				URI:            DefaultURI,
				ConnectTimeout: -1 * time.Second,
			},
			name:        "negative connect timeout",
			expectedErr: "connect timeout cannot be negative",
		},
		{
			cfg: &Config{
				URI:         DefaultURI,
				MaxPoolSize: testSmallPoolMax,
				MinPoolSize: testMinPoolLimit,
			},
			name:        "min pool size greater than max pool size",
			expectedErr: "min pool size cannot be greater than max pool size",
		},
		{
			cfg: &Config{
				URI:         DefaultURI,
				MaxPoolSize: testZeroPoolMax,
				MinPoolSize: testMinPoolLimit,
			},
			name:        "min pool size allowed when max pool size is zero",
			expectedErr: "",
		},
	}
}

// TestConfig_Validate tests validation rules on Config struct.
func TestConfig_Validate(t *testing.T) {
	for _, tt := range buildValidateCases() {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
