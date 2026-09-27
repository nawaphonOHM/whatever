package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testStandardPort = 27017
)

// validTestConfig creates a baseline valid Config instance.
func validTestConfig() *Config {
	cfg := DefaultConfig()
	cfg.Host = "localhost"
	cfg.Port = testStandardPort
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
	emptyHostCfg := validTestConfig()
	emptyHostCfg.Host = ""

	emptyUserCfg := validTestConfig()
	emptyUserCfg.Username = ""

	emptyPassCfg := validTestConfig()
	emptyPassCfg.Password = ""

	runValidateCases(t, []configValidateCase{
		{cfg: nil, name: "nil config", expectedErr: "mongodb config cannot be nil"},
		{cfg: emptyHostCfg, name: "empty Host", expectedErr: "mongodb host cannot be empty"},
		{cfg: emptyUserCfg, name: "empty Username", expectedErr: "mongodb username cannot be empty"},
		{cfg: emptyPassCfg, name: "empty Password", expectedErr: "mongodb password cannot be empty"},
		{cfg: validTestConfig(), name: "valid config", expectedErr: ""},
	})
}

// TestConfig_Validate_Protocols tests valid and invalid protocol values.
func TestConfig_Validate_Protocols(t *testing.T) {
	mongodbCfg := validTestConfig()
	mongodbCfg.Protocol = ProtocolMongoDB

	srvCfg := validTestConfig()
	srvCfg.Protocol = ProtocolMongoDBSrv

	invalidProtoCfg := validTestConfig()
	invalidProtoCfg.Protocol = "http"

	emptyProtoCfg := validTestConfig()
	emptyProtoCfg.Protocol = ""

	runValidateCases(t, []configValidateCase{
		{cfg: mongodbCfg, name: "valid mongodb protocol", expectedErr: ""},
		{cfg: srvCfg, name: "valid mongodb+srv protocol", expectedErr: ""},
		{cfg: invalidProtoCfg, name: "invalid http protocol", expectedErr: "invalid mongodb protocol: http"},
		{cfg: emptyProtoCfg, name: "empty protocol", expectedErr: "invalid mongodb protocol: "},
	})
}
