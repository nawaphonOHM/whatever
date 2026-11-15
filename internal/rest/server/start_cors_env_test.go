package server

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFromBluePrint_DefaultCORS verifies disabled CORS when unconfigured.
func TestNewFromBluePrint_DefaultCORS(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	setCORSConfigurationEnv(t, "", "", filepath.Join(testCORSFixtureDir, "valid_cors.yaml"))

	srv, err := NewFromBluePrint(createTestBlueprintWithRoute())
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin: testAppOrigin,
		method: http.MethodGet,
	})
	assert.NotEqual(t, http.StatusNoContent, preflight.Code)
	assertNoCORSHeaders(t, preflight)

	getResp := executeGET(srv, testPingPath, testAppOrigin)
	assert.Equal(t, http.StatusOK, getResp.Code)
	assertNoCORSHeaders(t, getResp)
}

// TestNewFromBluePrint_EnvironmentYAML verifies inline YAML configuration.
func TestNewFromBluePrint_EnvironmentYAML(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	jsonPath := filepath.Join(testCORSFixtureDir, "valid_cors.json")
	setCORSConfigurationEnv(t, readCORSFixture(t, "valid_cors.yaml"), "", jsonPath)

	srv, err := NewFromBluePrint(createTestBlueprintWithRoute())
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin: testAppOrigin,
		method: getMethod,
	})
	assert.Equal(t, http.StatusNoContent, preflight.Code)
	assertCORSHeaders(t, preflight, expectedHeaders{
		origin:  testAppOrigin,
		methods: yamlExpectedMethods,
		headers: "Authorization, Content-Type, Accept",
	})
}

// TestNewFromBluePrint_LegacyPathEnvironmentIsIgnored verifies legacy paths do not configure startup.
func TestNewFromBluePrint_LegacyPathEnvironmentIsIgnored(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	legacyEnvPath := filepath.Join(testCORSFixtureDir, "valid_cors.yaml")
	setCORSConfigurationEnv(t, "", "", legacyEnvPath)

	srv, err := NewFromBluePrint(createTestBlueprintWithRoute())
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin: testAppOrigin,
		method: getMethod,
	})
	assert.NotEqual(t, http.StatusNoContent, preflight.Code)
	assertNoCORSHeaders(t, preflight)
	assertDefaultCORSDisabledGET(t, srv)
}

func assertDefaultCORSDisabledGET(t *testing.T, srv *Server) {
	t.Helper()
	getResp := executeGET(srv, testPingPath, testAppOrigin)
	assert.Equal(t, http.StatusOK, getResp.Code)
	assertNoCORSHeaders(t, getResp)
}
