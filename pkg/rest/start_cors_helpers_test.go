package rest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
	envCORSPathKey      = "OHM9996_SERVER_CORS_PATH"
	envCORSYAMLKey      = "OHM9996_SERVER_CORS_YML_YAML"
	envCORSJSONKey      = "OHM9996_SERVER_CORS_JSON"
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

func newCORSBlueprintWithFixture() *BluePrint {
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
	return NewBluePrint().WithAPIs(regPing)
}

func testdataFixturePath(filename string) string {
	return filepath.Join("testdata", filename)
}

func readCORSTestFixture(t *testing.T, filename string) string {
	t.Helper()
	content, err := os.ReadFile(testdataFixturePath(filename))
	require.NoError(t, err)
	return string(content)
}
