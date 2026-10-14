package server

import (
	"testing"

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
