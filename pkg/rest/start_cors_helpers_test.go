package rest

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	e2ePingURLFormat    = "http://%s:%d/api/v1/ping"
	testAppOrigin       = "https://app.example.com"
	testAdminOrigin     = "https://admin.example.com"
	testHeaderOrigin    = "Origin"
	postMethod          = "POST"
	getMethod           = "GET"
	jsonExpectedMethods = "POST, GET"
	yamlExpectedMethods = "GET, POST"
	jsonExpectedHeaders = "Authorization, Content-Type, X-Request-ID"
	yamlExpectedHeaders = "Authorization, Content-Type, Accept"
	schemaValidationErr = "failed schema validation"
)

type preflightReq struct {
	origin  string
	method  string
	headers string
}

type expectedCORS struct {
	origin  string
	methods string
	headers string
}

func pingURL(port int) string {
	return fmt.Sprintf(e2ePingURLFormat, testHost, port)
}

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

func newCORSBlueprintWithFixture(corsPath string) *BluePrint {
	meta := NewMeta()
	if corsPath != "" {
		meta = meta.WithCorsPath(corsPath)
	}

	regPing := &RRestAPIRegistration{
		Version: 1,
		Prefix:  "/api",
		Apis: []*ExportableAPI{{
			Path:   "/ping",
			Method: GET,
			Handler: func(Context) Response {
				return OK("pong")
			},
		}},
	}
	return NewBluePrint().WithMeta(meta).WithAPIs(regPing)
}

func testdataFixturePath(filename string) string {
	return filepath.Join("testdata", filename)
}
