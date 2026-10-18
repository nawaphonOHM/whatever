package config

import (
	"testing"
)

const testErrInvalidEndpoint = "invalid endpoint"

// TestConfig_Validate_EndpointValid tests valid endpoint formats.
func TestConfig_Validate_EndpointValid(t *testing.T) {
	endpoints := []string{
		testDefaultEndpoint,
		"http://localhost:4317",
		"https://collector.internal:4318",
		"dns:///otel-collector:4317",
		"otel-collector",
		"127.0.0.1:4317",
		"[::1]:4317",
	}

	var cases []*configValidateCase
	for _, ep := range endpoints {
		cfg := validTestConfig()
		cfg.Endpoint = ep
		cases = append(cases, &configValidateCase{
			cfg:         cfg,
			name:        ep,
			expectedErr: "",
		})
	}
	runValidateCases(t, cases)
}

// createEndpointConfig returns a valid config with endpoint set.
func createEndpointConfig(endpoint string) *Config {
	cfg := validTestConfig()
	cfg.Endpoint = endpoint
	return cfg
}

// TestConfig_Validate_EndpointInvalid tests invalid endpoint strings.
func TestConfig_Validate_EndpointInvalid(t *testing.T) {
	runValidateCases(t, []*configValidateCase{
		{
			cfg:         createEndpointConfig(""),
			name:        "empty endpoint",
			expectedErr: "endpoint cannot be empty",
		},
		{
			cfg:         createEndpointConfig("localhost: 4317"),
			name:        "whitespace endpoint",
			expectedErr: testErrInvalidEndpoint,
		},
		{
			cfg:         createEndpointConfig("localhost:99999"),
			name:        "port above 65535",
			expectedErr: testErrInvalidEndpoint,
		},
		{
			cfg:         createEndpointConfig("localhost:0"),
			name:        "port zero",
			expectedErr: testErrInvalidEndpoint,
		},
		{
			cfg:         createEndpointConfig("localhost:invalid"),
			name:        "non-numeric port",
			expectedErr: testErrInvalidEndpoint,
		},
		{
			cfg:         createEndpointConfig(":4317"),
			name:        "empty host with port",
			expectedErr: testErrInvalidEndpoint,
		},
		{
			cfg:         createEndpointConfig("http://"),
			name:        "URL missing target",
			expectedErr: testErrInvalidEndpoint,
		},
		{
			cfg:         createEndpointConfig("http://localhost:99999"),
			name:        "URL port out of range",
			expectedErr: testErrInvalidEndpoint,
		},
	})
}
