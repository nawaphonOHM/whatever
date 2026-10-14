package server

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFromBluePrint_DefaultCORS verifies disabled CORS when unconfigured.
func TestNewFromBluePrint_DefaultCORS(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv("OHM9996_SERVER_CORS_PATH", "")

	srv, err := NewFromBluePrint(createTestBlueprintWithRoute(""))
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
