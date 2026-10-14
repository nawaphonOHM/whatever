package server

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const envCORSPathKey = "OHM9996_SERVER_CORS_PATH"

// TestNewFromBluePrint_EnvFallback verifies environment variable fallback when Meta.CorsPath is empty.
func TestNewFromBluePrint_EnvFallback(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envCORSPathKey, filepath.Join("testdata", "valid_cors.yaml"))

	srv, err := NewFromBluePrint(createTestBlueprintWithRoute(""))
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

// TestNewFromBluePrint_MetaPrecedenceOverEnv verifies Meta.CorsPath overrides environment variable.
func TestNewFromBluePrint_MetaPrecedenceOverEnv(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envCORSPathKey, filepath.Join("testdata", "valid_cors.yaml"))

	corsPath := filepath.Join("testdata", "valid_cors.json")
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute(corsPath))
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin: testAppOrigin,
		method: postMethod,
	})
	assert.Equal(t, http.StatusNoContent, preflight.Code)
	assertCORSHeaders(t, preflight, expectedHeaders{
		origin:  testAppOrigin,
		methods: jsonExpectedMethods,
		headers: jsonExpectedHeaders,
	})
}
