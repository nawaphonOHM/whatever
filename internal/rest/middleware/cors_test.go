// Package middleware provides unit tests for HTTP middlewares.
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
	testAPIData    = "/api/data"
	testOriginHdr  = "Origin"
)

// init initializes the gin test mode.
func init() {
	gin.SetMode(gin.TestMode)
}

// executePreflightRequest executes OPTIONS preflight request against router.
func executePreflightRequest(
	t *testing.T,
	r *gin.Engine,
	path, origin string,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, path, nil)
	require.NoError(t, err)
	req.Header.Set(testOriginHdr, origin)
	req.Header.Set("Access-Control-Request-Method", "GET")
	r.ServeHTTP(w, req)
	return w
}

// TestCORS_DefaultPreflight tests preflight OPTIONS request.
func TestCORS_DefaultPreflight(t *testing.T) {
	r := gin.New()
	r.Use(CORS(DefaultCORSConfig()))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := executePreflightRequest(t, r, "/test", "http://example.com")
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get(headerAllowOrg))
	assert.Contains(t, w.Header().Get(headerAllowMth), "GET")
}

// executeTestRequest runs a GET request with origin header against router.
func executeTestRequest(
	t *testing.T,
	r *gin.Engine,
	path, origin string,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)
	req.Header.Set(testOriginHdr, origin)
	r.ServeHTTP(w, req)
	return w
}

// setupCustomCORSEngine creates a test engine with custom CORS configuration.
func setupCustomCORSEngine() *gin.Engine {
	cfg := CORSConfig{
		AllowOrigins:     []string{testOriginApp},
		AllowMethods:     []string{"GET", "POST"},
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

// TestCORS_CustomConfig tests CORS middleware with customized origin policies.
func TestCORS_CustomConfig(t *testing.T) {
	r := setupCustomCORSEngine()

	// Matching origin request
	w := executeTestRequest(t, r, testAPIData, testOriginApp)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testOriginApp, w.Header().Get(headerAllowOrg))
	assert.Equal(t, "true", w.Header().Get(headerAllowCrd))
	assert.Equal(t, "X-Custom-Header", w.Header().Get(headerExposeHd))

	// Non-matching origin request
	w2 := executeTestRequest(t, r, testAPIData, testOriginEvil)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Empty(t, w2.Header().Get(headerAllowOrg))
}
