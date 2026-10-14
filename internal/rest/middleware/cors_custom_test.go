package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupCustomCORSEngine() *gin.Engine {
	cfg := &CORSConfig{
		AllowOrigins:     []string{testOriginApp},
		AllowMethods:     []string{methodGet, "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Custom-Header"},
		AllowCredentials: true,
		MaxAge:           "3600",
	}

	r := gin.New()
	r.Use(CORS(cfg))
	r.GET(testAPIData, func(c *gin.Context) {
		c.String(http.StatusOK, "data")
	})
	return r
}

func assertMatchingOriginHeaders(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, testOriginApp, w.Header().Get(headerAllowOrg))
	assert.Equal(t, headerValOrg, w.Header().Get(hdrVary))
	assert.Equal(t, "true", w.Header().Get(headerAllowCrd))
	assert.Equal(t, "X-Custom-Header", w.Header().Get(headerExposeHd))
}

func TestCORS_CustomConfig_MatchingOrigin(t *testing.T) {
	r := setupCustomCORSEngine()

	w := executeTestRequest(t, r, testAPIData, testOriginApp)
	assert.Equal(t, http.StatusOK, w.Code)
	assertMatchingOriginHeaders(t, w)

	wPreflight := executePreflightRequest(t, r, testAPIData, testOriginApp)
	assert.Equal(t, http.StatusNoContent, wPreflight.Code)
	assert.Equal(t, testOriginApp, wPreflight.Header().Get(headerAllowOrg))
	assert.Equal(t, headerValOrg, wPreflight.Header().Get(hdrVary))
}

func assertDisallowedHeaders(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Empty(t, w.Header().Get(headerAllowOrg))
	assert.Equal(t, headerValOrg, w.Header().Get(hdrVary))
}

func TestCORS_CustomConfig_DisallowedOrigin(t *testing.T) {
	r := setupCustomCORSEngine()

	w := executeTestRequest(t, r, testAPIData, testOriginEvil)
	assert.Equal(t, http.StatusOK, w.Code)
	assertDisallowedHeaders(t, w)

	wPre := executePreflightRequest(t, r, testAPIData, testOriginEvil)
	assert.NotEqual(t, http.StatusNoContent, wPre.Code)
	assertDisallowedHeaders(t, wPre)

	wNoOrigin := executeTestRequest(t, r, testAPIData, "")
	assert.Equal(t, http.StatusOK, wNoOrigin.Code)
	assertDisallowedHeaders(t, wNoOrigin)
}

func TestCORS_EmptySliceConfig(t *testing.T) {
	cfg := &CORSConfig{
		AllowOrigins: make([]string, 0),
		AllowMethods: []string{methodGet, "POST"},
	}
	r := gin.New()
	r.Use(CORS(cfg))
	r.GET(testPath, func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	wGet := executeTestRequest(t, r, testPath, testOriginEx)
	assert.Equal(t, http.StatusOK, wGet.Code)
	assertNoCORSResponse(t, wGet)
}
