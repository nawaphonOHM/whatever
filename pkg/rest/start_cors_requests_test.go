package rest

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setHeaderIfPresent(req *http.Request, key, val string) {
	if val != "" {
		req.Header.Set(key, val)
	}
}

func setPreflightHeaders(req *http.Request, opts preflightReq) {
	setHeaderIfPresent(req, "Origin", opts.origin)
	setHeaderIfPresent(req, "Access-Control-Request-Method", opts.method)
	setHeaderIfPresent(req, "Access-Control-Request-Headers", opts.headers)
}

func corsPreflightRequest(t *testing.T, url string, req preflightReq) *http.Response {
	t.Helper()
	httpReq, err := http.NewRequest(http.MethodOptions, url, nil)
	require.NoError(t, err)
	setPreflightHeaders(httpReq, req)

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	return resp
}

func corsGETRequest(t *testing.T, url, origin string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	setHeaderIfPresent(req, "Origin", origin)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func closeBody(t *testing.T, body io.ReadCloser) {
	t.Helper()
	_, err := io.Copy(io.Discard, body)
	assert.NoError(t, err)
	assert.NoError(t, body.Close())
	http.DefaultClient.CloseIdleConnections()
}
