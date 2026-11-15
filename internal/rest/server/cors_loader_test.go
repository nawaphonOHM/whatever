package server

import (
	"testing"

	intcfg "github.com/nawaphonOHM/whatever/v3/internal/rest/config"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCORSYAMLEnv             = "OHM9996_SERVER_CORS_YML_YAML"
	testCORSJSONEnv             = "OHM9996_SERVER_CORS_JSON"
	testCORSSentinel            = "https://private-cors-config.example.com"
	testYAMLConfig              = "allow-origins:\n  - https://yaml.example.com\n"
	testJSONConfig              = `{"allow-origins":["https://json.example.com"]}`
	testInvalidJSONOriginPrefix = `{"allow-origins":["`
	testInvalidJSONSchemaSuffix = `"],"unknown-property":true}`
)

func TestLoadCORSConfig_EnvironmentFormats(t *testing.T) {
	tests := []struct {
		name string
		cfg  HTTPFields
	}{
		{
			name: "yaml",
			cfg:  HTTPFields{CORSYAML: testYAMLConfig},
		},
		{
			name: "json",
			cfg:  HTTPFields{CORSJSON: testJSONConfig},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			corsCfg := loadCORSConfig(configWithHTTPFields(tc.cfg))
			assert.Equal(t, []string{"https://" + tc.name + ".example.com"}, corsCfg.AllowOrigins)
		})
	}
}

func TestLoadCORSConfig_JSONTakesPrecedence(t *testing.T) {
	cfg := configWithHTTPFields(HTTPFields{
		CORSYAML: testYAMLConfig,
		CORSJSON: testJSONConfig,
	})

	corsCfg := loadCORSConfig(cfg)
	assert.Equal(t, []string{"https://json.example.com"}, corsCfg.AllowOrigins)
}

func TestLoadCORSConfig_AbsentValuesLogInfoAndDisable(t *testing.T) {
	tests := []struct {
		cfg  *Config
		name string
	}{
		{name: "nil config"},
		{name: "empty config", cfg: &Config{}},
		{name: "whitespace values", cfg: configWithHTTPFields(HTTPFields{
			CORSYAML: " \t\n ",
			CORSJSON: " \t\n ",
		})},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf, reset := setupLifecycleLogger(t)
			defer reset()

			assert.Equal(t, middleware.DefaultCORSConfig(), loadCORSConfig(tc.cfg))
			entry := findLogEntry(parseJSONLogs(t, buf.Bytes()), msgCORSNotConfigured)
			require.NotNil(t, entry)
			assert.Equal(t, "INFO", entry[levelField])
		})
	}
}

func TestConfig_LoadsCORSEnvironmentValues(t *testing.T) {
	clearServerConfigEnv(t)
	t.Setenv(testCORSYAMLEnv, testYAMLConfig)
	t.Setenv(testCORSJSONEnv, testJSONConfig)

	cfg, err := intcfg.Load[Config]()
	require.NoError(t, err)
	assert.Equal(t, testYAMLConfig, cfg.CORSYAML)
	assert.Equal(t, testJSONConfig, cfg.CORSJSON)
}
