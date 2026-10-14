package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertJSONE2EPreflights(t *testing.T, url string) {
	appPre := corsPreflightRequest(t, url, preflightReq{
		origin:  testAppOrigin,
		method:  postMethod,
		headers: "Authorization, Content-Type",
	})
	defer closeBody(t, appPre.Body)
	assert.Equal(t, http.StatusNoContent, appPre.StatusCode)
	assertCORSHeaders(t, appPre, expectedCORS{
		origin:  testAppOrigin,
		methods: jsonExpectedMethods,
		headers: jsonExpectedHeaders,
	})

	disPre := corsPreflightRequest(t, url, preflightReq{
		origin: "https://malicious.org",
		method: postMethod,
	})
	defer closeBody(t, disPre.Body)
	assert.NotEqual(t, http.StatusNoContent, disPre.StatusCode)
	assertDisallowedCORS(t, disPre)
}

func assertJSONE2EGETs(t *testing.T, url string) {
	appGET := corsGETRequest(t, url, testAppOrigin)
	defer closeBody(t, appGET.Body)
	assert.Equal(t, http.StatusOK, appGET.StatusCode)
	assertCORSHeaders(t, appGET, expectedCORS{
		origin:  testAppOrigin,
		methods: jsonExpectedMethods,
		headers: jsonExpectedHeaders,
	})

	disGET := corsGETRequest(t, url, "https://malicious.org")
	defer closeBody(t, disGET.Body)
	assert.Equal(t, http.StatusOK, disGET.StatusCode)
	assertDisallowedCORS(t, disGET)
}

func TestStartREST_ValidJSONFixture(t *testing.T) {
	url, errCh := startCORSE2E(t, testdataFixturePath("valid_cors.json"))
	assertJSONE2EPreflights(t, url)
	assertJSONE2EGETs(t, url)
	signalAndWait(t, errCh)
}
