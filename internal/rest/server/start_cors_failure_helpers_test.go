package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readStartupCORSConfig(t *testing.T, tc serverFailureCase) (yamlConfig, jsonConfig string) {
	t.Helper()
	if tc.yamlFixture != "" {
		yamlConfig = readCORSFixture(t, tc.yamlFixture)
	}
	if tc.jsonFixture != "" {
		jsonConfig = readCORSFixture(t, tc.jsonFixture)
	}
	return yamlConfig, jsonConfig
}

func assertInvalidCORSDoesNotStopServer(t *testing.T) {
	t.Helper()
	srv, err := NewFromBluePrint(createTestBlueprintWithRoute())
	require.NoError(t, err)
	require.NotNil(t, srv)

	getResp := executeGET(srv, testPingPath, testAppOrigin)
	assert.Equal(t, http.StatusOK, getResp.Code)
	assertNoCORSHeaders(t, getResp)
	preflight := executePreflight(srv, testPingPath, preflightOpts{origin: testAppOrigin, method: getMethod})
	assert.NotEqual(t, http.StatusNoContent, preflight.Code)
	assertNoCORSHeaders(t, preflight)
}
