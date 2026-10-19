package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseYAMLCORS_MultiMethodRequests(t *testing.T) {
	yamlReqContent := `
allow-origins:
  - https://example.com
allow-methods:
  - POST
  - PUT
allow-headers:
  - Authorization
request-headers:
  - Content-Type
  - X-Request-ID
request-methods:
  - GET
  - DELETE
`
	_, fileCfg, err := parseYAMLCORS("cors_request.yaml", []byte(yamlReqContent))
	require.NoError(t, err)
	require.NotNil(t, fileCfg)
	assert.Equal(t, []string{"POST", "PUT"}, fileCfg.AllowMethods)
	assert.Equal(t, []string{"GET", "DELETE"}, fileCfg.RequestMethods)
	assert.Equal(t, []string{"Content-Type", "X-Request-ID"}, fileCfg.RequestHeaders)
}

func TestParseJSONCORS_MultiMethodRequests(t *testing.T) {
	jsonReqContent := `{
		"allow-origins": ["https://example.com"],
		"allow-methods": ["POST", "PUT"],
		"allow-headers": ["Authorization"],
		"request-headers": ["Content-Type", "X-Request-ID"],
		"request-methods": ["GET", "DELETE"]
	}`
	_, fileCfg, err := parseJSONCORS("cors_request.json", []byte(jsonReqContent))
	require.NoError(t, err)
	require.NotNil(t, fileCfg)
	assert.Equal(t, []string{"POST", "PUT"}, fileCfg.AllowMethods)
	assert.Equal(t, []string{"GET", "DELETE"}, fileCfg.RequestMethods)
	assert.Equal(t, []string{"Content-Type", "X-Request-ID"}, fileCfg.RequestHeaders)
}

func TestLoadCORSConfig_MetaOverridesServerConfig(t *testing.T) {
	metaPath := writeTempConfigFile(t, "meta_cors.json", `{"allow-origins": ["https://meta.example.com"]}`)
	envPath := writeTempConfigFile(t, "env_cors.json", `{"allow-origins": ["https://env.example.com"]}`)

	bp := contracts.NewBluePrint().WithMeta(contracts.NewMeta().WithCorsPath(metaPath))
	serverCfg := &Config{}
	serverCfg.CORSPath = envPath

	cfg, err := loadCORSConfig(serverCfg, bp)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://meta.example.com"}, cfg.AllowOrigins)
}

func TestLoadCORSConfig_FallbackToServerConfig(t *testing.T) {
	envPath := writeTempConfigFile(t, "env_cors.json", `{"allow-origins": ["https://env.example.com"]}`)

	bp := contracts.NewBluePrint().WithMeta(contracts.NewMeta())
	serverCfg := &Config{}
	serverCfg.CORSPath = envPath

	cfg, err := loadCORSConfig(serverCfg, bp)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://env.example.com"}, cfg.AllowOrigins)
}

func TestLoadCORSConfig_DefaultDisabledWhenNeither(t *testing.T) {
	bp := contracts.NewBluePrint().WithMeta(contracts.NewMeta())
	serverCfg := &Config{}

	cfg, err := loadCORSConfig(serverCfg, bp)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Nil(t, cfg.AllowOrigins)
}
