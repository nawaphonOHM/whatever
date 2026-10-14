package server

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFromBluePrint_ValidYAMLFixture verifies server initialization from YAML fixture.
func TestNewFromBluePrint_ValidYAMLFixture(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	corsPath := filepath.Join("testdata", "valid_cors.yaml")
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute(corsPath))
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin:  testAdminOrigin,
		method:  getMethod,
		headers: "Authorization",
	})
	assert.Equal(t, http.StatusNoContent, preflight.Code)
	assertCORSHeaders(t, preflight, expectedHeaders{
		origin:  testAdminOrigin,
		methods: getMethod,
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
	corsPath := filepath.Join("testdata", "valid_cors_with_requests.yaml")
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute(corsPath))
	require.NoError(t, err)

	preflight := executePreflight(srv, testPingPath, preflightOpts{
		origin:  testAppOrigin,
		method:  http.MethodPut,
		headers: "Authorization",
	})
	assert.Equal(t, http.StatusNoContent, preflight.Code)
	assertCORSHeaders(t, preflight, expectedHeaders{
		origin:  testAppOrigin,
		methods: "PUT",
		headers: "Authorization, Content-Type",
	})
}
