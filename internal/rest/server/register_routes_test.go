package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/health"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIVersion = "1.2.3"
	itemsPath      = "/v1/items"
)

// sampleRegs returns a simple versioned registration set.
func sampleRegs() []*RestAPIRegistration {
	return []*RestAPIRegistration{{
		Version: 1,
		Prefix:  "/items",
		Apis: []*ExportableAPI{{
			Path:   "",
			Method: GET,
			Handler: func(*Context) response.Response {
				return response.OK([]string{"item1", "item2"})
			},
		}},
	}}
}

// getJSON performs a GET and returns status plus body bytes.
func getJSON(
	t *testing.T,
	engine *gin.Engine,
	path string,
) (int, []byte) {
	t.Helper()
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w.Code, w.Body.Bytes()
}

// assertItemsOK checks the versioned items payload.
func assertItemsOK(t *testing.T, body []byte) {
	t.Helper()
	var itemsResp response.SuccessResponse[[]string]
	require.NoError(t, json.Unmarshal(body, &itemsResp))
	assert.True(t, itemsResp.Success)
	assert.Equal(t, []string{"item1", "item2"}, itemsResp.Data)
}

// assertHealthVersion checks reserved health payload version.
func assertHealthVersion(t *testing.T, body []byte) {
	t.Helper()
	var healthResp response.SuccessResponse[health.Status]
	require.NoError(t, json.Unmarshal(body, &healthResp))
	assert.Equal(t, testAPIVersion, healthResp.Data.Version)
}

// TestRegisterRoutes_Success mounts business and health routes.
func TestRegisterRoutes_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	require.NoError(t, RegisterRoutesWithVersion(
		engine, sampleRegs(), testAPIVersion,
	))

	// Act / Assert business route
	code, body := getJSON(t, engine, itemsPath)
	assert.Equal(t, http.StatusOK, code)
	assertItemsOK(t, body)

	// Act / Assert reserved health
	code, body = getJSON(t, engine, ReservedHealthPath)
	assert.Equal(t, http.StatusOK, code)
	assertHealthVersion(t, body)
}
