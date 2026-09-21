package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validTestConfig creates a baseline valid Config instance.
func validTestConfig() *Config {
	cfg := DefaultConfig()
	cfg.URI = "mongodb://localhost:27017"
	cfg.Database = "testdb"
	cfg.Username = "testuser"
	cfg.Password = "testpassword"
	return cfg
}

// configValidateCase defines a test case for Config.Validate.
type configValidateCase struct {
	cfg         *Config
	name        string
	expectedErr string
}

func runValidateCases(t *testing.T, cases []configValidateCase) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

// TestConfig_Validate_RequiredFields tests required field validation rules on Config struct.
func TestConfig_Validate_RequiredFields(t *testing.T) {
	emptyURICfg := validTestConfig()
	emptyURICfg.URI = ""

	emptyDBCfg := validTestConfig()
	emptyDBCfg.Database = ""

	emptyUserCfg := validTestConfig()
	emptyUserCfg.Username = ""

	emptyPassCfg := validTestConfig()
	emptyPassCfg.Password = ""

	runValidateCases(t, []configValidateCase{
		{cfg: nil, name: "nil config", expectedErr: "mongodb config cannot be nil"},
		{cfg: emptyURICfg, name: "empty URI", expectedErr: "mongodb uri cannot be empty"},
		{cfg: emptyDBCfg, name: "empty Database", expectedErr: "mongodb database cannot be empty"},
		{cfg: emptyUserCfg, name: "empty Username", expectedErr: "mongodb username cannot be empty"},
		{cfg: emptyPassCfg, name: "empty Password", expectedErr: "mongodb password cannot be empty"},
		{cfg: validTestConfig(), name: "valid config", expectedErr: ""},
	})
}
