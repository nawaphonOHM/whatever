package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/stretchr/testify/require"
)

const (
	testFilePerm        = 0o600
	testPingPath        = "/api/v1/ping"
	testAppOrigin       = "https://app.example.com"
	testAdminOrigin     = "https://admin.example.com"
	originHeaderVal     = "Origin"
	patternMismatch     = "Does not match pattern"
	schemaValidationErr = "failed schema validation"
	postMethod          = "POST"
	getMethod           = "GET"
)

type preflightOpts struct {
	origin  string
	method  string
	headers string
}

type expectedHeaders struct {
	origin  string
	methods string
	headers string
}

func writeTempConfigFile(t *testing.T, filename, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), testFilePerm)
	require.NoError(t, err)
	return path
}

func createTestBlueprintWithRoute(corsPath string) *contracts.BluePrint {
	meta := contracts.NewMeta()
	if corsPath != "" {
		meta = meta.WithCorsPath(corsPath)
	}

	api := &contracts.ExportableAPI{
		Path:   "/ping",
		Method: contracts.GET,
		Handler: func(contracts.Context) contracts.Response {
			return testOK("pong")
		},
	}
	reg := &contracts.RRestAPIRegistration{
		Version: 1,
		Prefix:  "/api",
		Apis:    []*contracts.ExportableAPI{api},
	}
	return contracts.NewBluePrint().WithMeta(meta).WithAPIs(reg)
}

func setReqHeader(req *http.Request, key, val string) {
	if val != "" {
		req.Header.Set(key, val)
	}
}

func executePreflight(srv *Server, path string, opts preflightOpts) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, path, nil)
	setReqHeader(req, "Origin", opts.origin)
	setReqHeader(req, "Access-Control-Request-Method", opts.method)
	setReqHeader(req, "Access-Control-Request-Headers", opts.headers)
	srv.Engine.ServeHTTP(w, req)
	return w
}

func executeGET(srv *Server, path, origin string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	setReqHeader(req, "Origin", origin)
	srv.Engine.ServeHTTP(w, req)
	return w
}
