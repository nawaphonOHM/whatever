package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func startCORSE2EWithLegacyPath(t *testing.T, envPath string) (string, <-chan error) {
	t.Helper()
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envCORSYAMLKey, "")
	t.Setenv(envCORSJSONKey, "")
	t.Setenv(envCORSPathKey, envPath)
	errCh := startRESTAsync(newCORSBlueprintWithFixture())
	return pingURL(port), errCh
}

func assertEnvYAMLPreflight(t *testing.T, url string) {
	t.Helper()
	preflight := corsPreflightRequest(t, url, preflightReq{
		origin: testAppOrigin,
		method: getMethod,
	})
	defer closeBody(t, preflight.Body)
	assert.Equal(t, http.StatusNoContent, preflight.StatusCode)
	assertCORSHeaders(t, preflight, expectedCORS{
		origin:  testAppOrigin,
		methods: yamlExpectedMethods,
		headers: yamlExpectedHeaders,
	})
}

func assertEnvJSONPreflight(t *testing.T, url string) {
	t.Helper()
	preflight := corsPreflightRequest(t, url, preflightReq{
		origin: testAppOrigin,
		method: postMethod,
	})
	defer closeBody(t, preflight.Body)
	assert.Equal(t, http.StatusNoContent, preflight.StatusCode)
	assertCORSHeaders(t, preflight, expectedCORS{
		origin:  testAppOrigin,
		methods: jsonExpectedMethods,
		headers: jsonExpectedHeaders,
	})
}

func TestStartREST_EnvironmentYAML(t *testing.T) {
	url, errCh := startCORSE2E(t, readCORSTestFixture(t, "valid_cors.yaml"), "")
	assertEnvYAMLPreflight(t, url)
	signalAndWait(t, errCh)
}

func TestStartREST_JSONPrecedence(t *testing.T) {
	url, errCh := startCORSE2E(t, readCORSTestFixture(t, "valid_cors.yaml"), readCORSTestFixture(t, "valid_cors.json"))
	assertEnvJSONPreflight(t, url)
	signalAndWait(t, errCh)
}

func TestStartREST_LegacyPathsAreIgnored(t *testing.T) {
	envPath := testdataFixturePath("valid_cors.yaml")
	url, errCh := startCORSE2EWithLegacyPath(t, envPath)
	assertDefaultDisabledPreflight(t, url)
	assertDefaultDisabledGET(t, url)
	signalAndWait(t, errCh)
}
