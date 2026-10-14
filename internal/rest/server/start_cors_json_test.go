package server

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jsonExpectedHeaders = "Authorization, Content-Type, X-Request-ID"

func assertJSONPreflights(t *testing.T, srv *Server) {
	appPre := executePreflight(srv, testPingPath, preflightOpts{
		origin:  testAppOrigin,
		method:  postMethod,
		headers: "Authorization, Content-Type",
	})
	assert.Equal(t, http.StatusNoContent, appPre.Code)
	assertCORSHeaders(t, appPre, expectedHeaders{
		origin:  testAppOrigin,
		methods: postMethod,
		headers: jsonExpectedHeaders,
	})

	disPre := executePreflight(srv, testPingPath, preflightOpts{
		origin: "https://malicious.org",
		method: postMethod,
	})
	assert.NotEqual(t, http.StatusNoContent, disPre.Code)
	assertDisallowedCORS(t, disPre)
}

func assertJSONGETs(t *testing.T, srv *Server) {
	appGET := executeGET(srv, testPingPath, testAppOrigin)
	assert.Equal(t, http.StatusOK, appGET.Code)
	assertCORSHeaders(t, appGET, expectedHeaders{
		origin:  testAppOrigin,
		methods: postMethod,
		headers: jsonExpectedHeaders,
	})

	disGET := executeGET(srv, testPingPath, "https://malicious.org")
	assert.Equal(t, http.StatusOK, disGET.Code)
	assertDisallowedCORS(t, disGET)

	noOrgGET := executeGET(srv, testPingPath, "")
	assert.Equal(t, http.StatusOK, noOrgGET.Code)
	assertDisallowedCORS(t, noOrgGET)
}

// TestNewFromBluePrint_ValidJSONFixture verifies server initialization from JSON fixture.
func TestNewFromBluePrint_ValidJSONFixture(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	corsPath := filepath.Join("testdata", "valid_cors.json")
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute(corsPath))
	require.NoError(t, err)

	assertJSONPreflights(t, srv)
	assertJSONGETs(t, srv)
}
