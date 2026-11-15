package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configWithHTTPFields(fields HTTPFields) *Config {
	cfg := &Config{}
	cfg.HTTPFields = fields
	return cfg
}

func TestLoadCORSConfig_InvalidPolicyLogsWarnAndDoesNotExposeContent(t *testing.T) {
	tests := []struct {
		name string
		cfg  HTTPFields
	}{
		{
			name: "invalid JSON syntax",
			cfg:  HTTPFields{CORSJSON: testInvalidJSONOriginPrefix + testCORSSentinel + `"`},
		},
		{
			name: "invalid JSON schema",
			cfg:  HTTPFields{CORSJSON: testInvalidJSONOriginPrefix + testCORSSentinel + testInvalidJSONSchemaSuffix},
		},
		{
			name: "trailing JSON document",
			cfg:  HTTPFields{CORSJSON: testInvalidJSONOriginPrefix + testCORSSentinel + `"]} {}`},
		},
		{
			name: "invalid YAML syntax",
			cfg:  HTTPFields{CORSYAML: "allow-origins: [\n  - " + testCORSSentinel},
		},
		{
			name: "trailing YAML document",
			cfg:  HTTPFields{CORSYAML: "allow-origins:\n  - " + testCORSSentinel + "\n---\nallow-origins: []\n"},
		},
		{
			name: "invalid JSON does not fall back to valid YAML",
			cfg: HTTPFields{
				CORSYAML: testYAMLConfig,
				CORSJSON: testInvalidJSONOriginPrefix + testCORSSentinel + testInvalidJSONSchemaSuffix,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf, reset := setupLifecycleLogger(t)
			defer reset()

			corsCfg := loadCORSConfig(configWithHTTPFields(tc.cfg))
			assert.Equal(t, middleware.DefaultCORSConfig(), corsCfg)
			entry := findLogEntry(parseJSONLogs(t, buf.Bytes()), msgCORSInvalid)
			require.NotNil(t, entry)
			assert.Equal(t, "WARN", entry[levelField])
			assert.NotContains(t, buf.String(), testCORSSentinel)
		})
	}
}
