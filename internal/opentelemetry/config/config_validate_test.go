package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// configValidateCase defines a test case for Config.Validate.
type configValidateCase struct {
	cfg         *Config
	name        string
	expectedErr string
}

// validTestConfig creates a baseline valid Config instance.
func validTestConfig() *Config {
	cfg := DefaultConfig()
	cfg.ServiceName = testServiceName
	cfg.Endpoint = testDefaultEndpoint
	cfg.Protocol = ProtocolGRPC
	cfg.SampleRate = DefaultSampleRate
	cfg.Enabled = true
	cfg.Insecure = true
	return cfg
}

// runValidateCases executes table-driven validation test cases.
func runValidateCases(t *testing.T, cases []configValidateCase) {
	t.Helper()
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

// TestConfig_Validate_NilAndValid tests nil and valid baseline configurations.
func TestConfig_Validate_NilAndValid(t *testing.T) {
	disabledCfg := &Config{
		Enabled:    false,
		SampleRate: DefaultSampleRate,
	}
	runValidateCases(t, []configValidateCase{
		{cfg: nil, name: "nil config", expectedErr: "opentelemetry config cannot be nil"},
		{cfg: validTestConfig(), name: "valid default config", expectedErr: ""},
		{cfg: disabledCfg, name: "valid disabled config with empty fields", expectedErr: ""},
	})
}

// TestConfig_Validate_ServiceName tests service name validation.
func TestConfig_Validate_ServiceName(t *testing.T) {
	emptyNameCfg := validTestConfig()
	emptyNameCfg.ServiceName = ""

	whitespaceNameCfg := validTestConfig()
	whitespaceNameCfg.ServiceName = "   "

	runValidateCases(t, []configValidateCase{
		{cfg: emptyNameCfg, name: "empty service name", expectedErr: "service name cannot be empty"},
		{cfg: whitespaceNameCfg, name: "whitespace service name", expectedErr: "service name cannot be empty"},
	})
}

// TestConfig_Validate_Protocols tests valid and invalid protocol strings.
func TestConfig_Validate_Protocols(t *testing.T) {
	protoCfg := validTestConfig()
	protoCfg.Protocol = ProtocolHTTPProtobuf

	jsonCfg := validTestConfig()
	jsonCfg.Protocol = ProtocolHTTPJSON

	invalidCfg := validTestConfig()
	invalidCfg.Protocol = "invalid"

	runValidateCases(t, []configValidateCase{
		{cfg: validTestConfig(), name: "valid grpc protocol", expectedErr: ""},
		{cfg: protoCfg, name: "valid http/protobuf protocol", expectedErr: ""},
		{cfg: jsonCfg, name: "valid http/json protocol", expectedErr: ""},
		{cfg: invalidCfg, name: "invalid protocol", expectedErr: "invalid protocol"},
	})
}
