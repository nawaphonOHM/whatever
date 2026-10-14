package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const envCORSKey = "OHM9996_SERVER_CORS_PATH"

func startCORSE2EWithEnv(t *testing.T, metaPath, envPath string) (string, <-chan error) {
	t.Helper()
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envCORSKey, envPath)
	errCh := startRESTAsync(newCORSBlueprintWithFixture(metaPath))
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
		methods: getMethod,
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
		methods: postMethod,
		headers: jsonExpectedHeaders,
	})
}

func TestStartREST_EnvFallback(t *testing.T) {
	url, errCh := startCORSE2EWithEnv(t, "", testdataFixturePath("valid_cors.yaml"))
	assertEnvYAMLPreflight(t, url)
	signalAndWait(t, errCh)
}

func TestStartREST_MetaPrecedenceOverEnv(t *testing.T) {
	jsonPath := testdataFixturePath("valid_cors.json")
	yamlPath := testdataFixturePath("valid_cors.yaml")
	url, errCh := startCORSE2EWithEnv(t, jsonPath, yamlPath)
	assertEnvJSONPreflight(t, url)
	signalAndWait(t, errCh)
}
