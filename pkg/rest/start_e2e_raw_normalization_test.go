package rest

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRawNormalizationBlueprint() *BluePrint {
	regMultiSlash := &RRestAPIRegistration{
		Group: APIGroupRaw, Prefix: "///raw-prefix///",
		Apis: []*ExportableAPI{
			{Path: "///endpoint///", Method: GET, Handler: func(Context) Response { return OK("multi-slash") }},
			{Path: "/item", Method: POST, Handler: func(Context) Response { return Created("item-created") }},
			{Path: "/item", Method: PUT, Handler: func(Context) Response { return OK("item-updated") }},
			{Path: "/item", Method: DELETE, Handler: func(Context) Response { return NoContent() }},
		},
	}
	return NewBluePrint().WithAPIs(regMultiSlash)
}

func assertRawGetAndPost(t *testing.T, baseURL string) {
	t.Helper()
	status, body := httpGetBody(t, baseURL+"/raw-prefix/endpoint")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "multi-slash")

	postStatus, postBody := httpPostBody(t, baseURL+"/raw-prefix/item", "{}")
	assert.Equal(t, http.StatusCreated, postStatus)
	assert.Contains(t, postBody, "item-created")
}

func assertRawPut(t *testing.T, baseURL string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPut, baseURL+"/raw-prefix/item", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		t.Logf("failed to discard body: %v", copyErr)
	}
	assert.NoError(t, resp.Body.Close())
	http.DefaultClient.CloseIdleConnections()
}

func assertRawDelete(t *testing.T, baseURL string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/raw-prefix/item", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		t.Logf("failed to discard body: %v", copyErr)
	}
	assert.NoError(t, resp.Body.Close())
	http.DefaultClient.CloseIdleConnections()
}

func TestStartREST_RawNormalizationAndMethods(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, testGroupRaw)

	errCh := startRESTAsync(newRawNormalizationBlueprint())
	baseURL := fmt.Sprintf("http://%s:%d", testHost, port)

	assertRawGetAndPost(t, baseURL)
	assertRawPut(t, baseURL)
	assertRawDelete(t, baseURL)
	signalAndWait(t, errCh)
}
