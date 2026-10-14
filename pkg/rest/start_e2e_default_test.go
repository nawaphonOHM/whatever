package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertDefaultDisabledPreflight(t *testing.T, url string) {
	t.Helper()
	preflight := corsPreflightRequest(t, url, preflightReq{
		origin: "https://example.com",
		method: http.MethodGet,
	})
	defer closeBody(t, preflight.Body)
	assert.NotEqual(t, http.StatusNoContent, preflight.StatusCode)
	assertNoCORSHeaders(t, preflight)
}

func assertDefaultDisabledGET(t *testing.T, url string) {
	t.Helper()
	getResp := corsGETRequest(t, url, "https://example.com")
	defer closeBody(t, getResp.Body)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	assertNoCORSHeaders(t, getResp)
}

func TestStartREST_DefaultCORS(t *testing.T) {
	url, errCh := startCORSE2E(t, "")
	assertDefaultDisabledPreflight(t, url)
	assertDefaultDisabledGET(t, url)
	signalAndWait(t, errCh)
}
