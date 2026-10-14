package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertYAMLE2EAllowed(t *testing.T, url string) {
	t.Helper()
	preflight := corsPreflightRequest(t, url, preflightReq{
		origin:  testAdminOrigin,
		method:  getMethod,
		headers: "Authorization",
	})
	defer closeBody(t, preflight.Body)
	assert.Equal(t, http.StatusNoContent, preflight.StatusCode)
	assertCORSHeaders(t, preflight, expectedCORS{
		origin:  testAdminOrigin,
		methods: getMethod,
		headers: yamlExpectedHeaders,
	})
}

func assertYAMLE2EDisallowed(t *testing.T, url string) {
	t.Helper()
	disPre := corsPreflightRequest(t, url, preflightReq{
		origin: "https://unauthorized.org",
		method: getMethod,
	})
	defer closeBody(t, disPre.Body)
	assert.NotEqual(t, http.StatusNoContent, disPre.StatusCode)
	assertDisallowedCORS(t, disPre)
}

func TestStartREST_ValidYAMLFixture(t *testing.T) {
	url, errCh := startCORSE2E(t, testdataFixturePath("valid_cors.yaml"))
	assertYAMLE2EAllowed(t, url)
	assertYAMLE2EDisallowed(t, url)
	signalAndWait(t, errCh)
}

func TestStartREST_ValidYAMLWithRequestsFixture(t *testing.T) {
	path := testdataFixturePath("valid_cors_with_requests.yaml")
	url, errCh := startCORSE2E(t, path)
	assertYAMLRequestsPreflight(t, url)
	signalAndWait(t, errCh)
}

func assertYAMLRequestsPreflight(t *testing.T, url string) {
	t.Helper()
	preflight := corsPreflightRequest(t, url, preflightReq{
		origin:  testAppOrigin,
		method:  http.MethodPut,
		headers: "Authorization",
	})
	defer closeBody(t, preflight.Body)
	assert.Equal(t, http.StatusNoContent, preflight.StatusCode)
	assertCORSHeaders(t, preflight, expectedCORS{
		origin:  testAppOrigin,
		methods: "PUT",
		headers: "Authorization, Content-Type",
	})
}
