package server

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFromBluePrint_ValidYAMLFixture verifies server initialization from YAML fixture.
func TestNewFromBluePrint_ValidYAMLFixture(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	setCORSConfigurationEnv(t, readCORSFixture(t, "valid_cors.yaml"), "", "")
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute())
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin:  testAdminOrigin,
		method:  getMethod,
		headers: "Authorization",
	})
	assert.Equal(t, http.StatusNoContent, preflight.Code)
	assertCORSHeaders(t, preflight, expectedHeaders{
		origin:  testAdminOrigin,
		methods: yamlExpectedMethods,
		headers: "Authorization, Content-Type, Accept",
	})

	disPreflight := executePreflight(srv, testPingPath, preflightOpts{
		origin: "https://unauthorized.org",
		method: getMethod,
	})
	assert.NotEqual(t, http.StatusNoContent, disPreflight.Code)
	assertDisallowedCORS(t, disPreflight)
}

// TestNewFromBluePrint_ValidYAMLWithRequestsFixture verifies server initialization with request fields.
func TestNewFromBluePrint_ValidYAMLWithRequestsFixture(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	setCORSConfigurationEnv(t, readCORSFixture(t, "valid_cors_with_requests.yaml"), "", "")
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute())
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin:  testAppOrigin,
		method:  http.MethodPut,
		headers: "Authorization",
	})
	assert.Equal(t, http.StatusNoContent, preflight.Code)
	assertCORSHeaders(t, preflight, expectedHeaders{
		origin:  testAppOrigin,
		methods: "PUT, DELETE",
		headers: "Authorization, Content-Type",
	})
}

type serverFailureCase struct {
	name        string
	yamlFixture string
	jsonFixture string
}

func getStartupFailureCases() []serverFailureCase {
	return []serverFailureCase{
		{name: "invalid JSON syntax", jsonFixture: "invalid_syntax.json"},
		{name: "invalid YAML syntax", yamlFixture: "invalid_syntax.yaml"},
		{name: "multi-document YAML", yamlFixture: "invalid_multidoc.yaml"},
		{name: "disallowed method", jsonFixture: "invalid_schema_disallowed_method.json"},
		{name: "bad origin format", jsonFixture: "invalid_schema_bad_origin.json"},
		{name: "bad header format", jsonFixture: "invalid_schema_bad_header.json"},
		{name: "unknown property", jsonFixture: "invalid_unknown_property.json"},
		{name: "YAML disallowed method", yamlFixture: "invalid_schema_disallowed_method.yaml"},
		{name: "YAML bad origin format", yamlFixture: "invalid_schema_bad_origin.yaml"},
		{
			name:        "invalid selected JSON does not fall back to valid YAML",
			yamlFixture: "valid_cors.yaml",
			jsonFixture: "invalid_unknown_property.json",
		},
	}
}

// TestNewFromBluePrint_StartupFailures verifies invalid CORS does not block server startup.
func TestNewFromBluePrint_StartupFailures(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)

	for _, tc := range getStartupFailureCases() {
		t.Run(tc.name, func(t *testing.T) {
			yamlConfig, jsonConfig := readStartupCORSConfig(t, tc)
			setCORSConfigurationEnv(t, yamlConfig, jsonConfig, "")
			assertInvalidCORSDoesNotStopServer(t)
		})
	}
}
