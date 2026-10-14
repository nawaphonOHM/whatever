package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAndValidateCORS_EmptyPath(t *testing.T) {
	cfg, err := LoadAndValidateCORS("")
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Nil(t, cfg.AllowOrigins)
	assert.Empty(t, cfg.AllowMethods)
	assert.Empty(t, cfg.AllowHeaders)
}

func TestLoadAndValidateCORS_ValidJSON(t *testing.T) {
	content := `{
		"allow-origins": ["https://example.com", "http://localhost.dev"],
		"allow-methods": "POST",
		"allow-headers": ["Authorization", "Content-Type", "Accept"]
	}`
	path := writeTempConfigFile(t, "cors.json", content)

	cfg, err := LoadAndValidateCORS(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, []string{"https://example.com", "http://localhost.dev"}, cfg.AllowOrigins)
	assert.Equal(t, []string{"POST"}, cfg.AllowMethods)
	assert.Equal(t, []string{"Authorization", "Content-Type", "Accept"}, cfg.AllowHeaders)
}

func TestLoadAndValidateCORS_ValidYAML(t *testing.T) {
	yamlContent := `
allow-origins:
  - https://api.example.com
  - https://app.example.com/v1
allow-methods: GET
allow-headers:
  - Content-Type
  - X-Request-ID
`
	path := writeTempConfigFile(t, "cors.yaml", yamlContent)
	cfg, err := LoadAndValidateCORS(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, []string{"https://api.example.com", "https://app.example.com/v1"}, cfg.AllowOrigins)
	assert.Equal(t, []string{"GET"}, cfg.AllowMethods)
	assert.Equal(t, []string{"Content-Type", "X-Request-ID"}, cfg.AllowHeaders)
}

func TestLoadAndValidateCORS_ValidYAMLWithRequests(t *testing.T) {
	yamlReqContent := `
allow-origins:
  - https://example.com
allow-methods: POST
allow-headers:
  - Authorization
request-headers:
  - Content-Type
request-methods: GET
`
	path := writeTempConfigFile(t, "cors_request.yaml", yamlReqContent)
	cfg, err := LoadAndValidateCORS(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, []string{"https://example.com"}, cfg.AllowOrigins)
	assert.Equal(t, []string{"POST"}, cfg.AllowMethods)
	assert.Equal(t, []string{"Authorization"}, cfg.AllowHeaders)
}
