package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testOriginApp  = "https://app.example.com"
	testOriginEvil = "https://evil.com"
	testOriginEx   = "http://example.com"
	testAPIData    = "/api/data"
	testOriginHdr  = "Origin"
	testPath       = "/test"
	methodGet      = "GET"
	hdrVary        = "Vary"
)

func init() { gin.SetMode(gin.TestMode) }

func executePreflightRequest(t *testing.T, r *gin.Engine, path, origin string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, path, nil)
	require.NoError(t, err)
	req.Header.Set(testOriginHdr, origin)
	req.Header.Set("Access-Control-Request-Method", methodGet)
	r.ServeHTTP(w, req)
	return w
}

func executeTestRequest(t *testing.T, r *gin.Engine, path, origin string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)
	if origin != "" {
		req.Header.Set(testOriginHdr, origin)
	}
	r.ServeHTTP(w, req)
	return w
}

func assertNoCORSResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Empty(t, w.Header().Get(headerAllowOrg))
	assert.Empty(t, w.Header().Get(hdrVary))
}

func TestCORS_DefaultDisabled(t *testing.T) {
	r := gin.New()
	r.Use(CORS(DefaultCORSConfig()))
	r.GET(testPath, func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	wPreflight := executePreflightRequest(t, r, testPath, testOriginEx)
	assert.NotEqual(t, http.StatusNoContent, wPreflight.Code)
	assertNoCORSResponse(t, wPreflight)
	wGet := executeTestRequest(t, r, testPath, testOriginEx)
	assert.Equal(t, http.StatusOK, wGet.Code)
	assertNoCORSResponse(t, wGet)
}

func TestCORS_NilAndEmptyConfig(t *testing.T) {
	r := gin.New()
	r.Use(CORS(nil))
	r.GET(testPath, func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	wGet := executeTestRequest(t, r, testPath, testOriginEx)
	assert.Equal(t, http.StatusOK, wGet.Code)
	assertNoCORSResponse(t, wGet)
}

func assertWildcardResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, wildcardOrigin, w.Header().Get(headerAllowOrg))
	assert.Empty(t, w.Header().Get(hdrVary))
}

func TestCORS_WildcardConfig(t *testing.T) {
	cfg := &CORSConfig{
		AllowOrigins: []string{wildcardOrigin},
		AllowMethods: []string{methodGet, "POST"},
	}
	r := gin.New()
	r.Use(CORS(cfg))
	r.GET(testPath, func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	wPre := executePreflightRequest(t, r, testPath, testOriginEx)
	assert.Equal(t, http.StatusNoContent, wPre.Code)
	assertWildcardResponse(t, wPre)
	wGet := executeTestRequest(t, r, testPath, testOriginEx)
	assert.Equal(t, http.StatusOK, wGet.Code)
	assertWildcardResponse(t, wGet)
}
