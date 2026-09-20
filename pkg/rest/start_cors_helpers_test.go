package rest

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// corsPreflightRequest sends an OPTIONS preflight request to the given URL.
func corsPreflightRequest(
	t *testing.T,
	url string,
) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

// closeBody closes the response body and asserts no error.
func closeBody(t *testing.T, body io.Closer) {
	t.Helper()
	assert.NoError(t, body.Close())
}

// newCustomCORSBlueprint creates a blueprint with custom CORS settings.
func newCustomCORSBlueprint() *BluePrint {
	cors := NewCorsSetting().
		WithAllowOrigin("https://example.com").
		WithAllowHTTPMethods(POST, PATCH)
	meta := NewMeta().WithCors(cors)

	api := &ExportableAPI{
		Path:   "/test",
		Method: GET,
		Handler: func(Context) Response {
			return OK("ok")
		},
	}
	reg := &RRestAPIRegistration{
		Apis: []*ExportableAPI{api},
	}
	return NewBluePrint().WithMeta(meta).WithAPIs(reg)
}

// assertCORSHeaders verifies CORS response headers.
func assertCORSHeaders(
	t *testing.T,
	resp *http.Response,
	expectedOrigin,
	expectedMethods string,
) {
	t.Helper()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, expectedOrigin,
		resp.Header.Get("Access-Control-Allow-Origin"))
	if expectedMethods != "" {
		assert.Equal(t, expectedMethods,
			resp.Header.Get("Access-Control-Allow-Methods"))
	} else {
		assert.NotEmpty(t,
			resp.Header.Get("Access-Control-Allow-Methods"))
	}
}
