package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertCORSHeaders(t *testing.T, resp *http.Response, exp expectedCORS) {
	t.Helper()
	assert.Equal(t, exp.origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, testHeaderOrigin, resp.Header.Get("Vary"))
	assert.Equal(t, exp.methods, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, exp.headers, resp.Header.Get("Access-Control-Allow-Headers"))
}

func assertNoCORSHeaders(t *testing.T, resp *http.Response) {
	t.Helper()
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Vary"))
}

func assertDisallowedCORS(t *testing.T, resp *http.Response) {
	t.Helper()
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, testHeaderOrigin, resp.Header.Get("Vary"))
}
