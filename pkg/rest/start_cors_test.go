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

func startCORSE2E(t *testing.T, corsPath string) (string, <-chan error) {
	t.Helper()
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(newCORSBlueprintWithFixture(corsPath))
	return pingURL(port), errCh
}

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
