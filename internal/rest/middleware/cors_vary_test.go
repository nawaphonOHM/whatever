package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runVaryRequest(t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, testPath, nil)
	require.NoError(t, err)
	req.Header.Set(testOriginHdr, testOriginApp)
	r.ServeHTTP(w, req)
	return w
}

func setupVaryEngineMultiple() *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Add(hdrVary, "Accept-Encoding")
		c.Writer.Header().Add(hdrVary, "X-Origin-Hint")
		c.Next()
	})
	r.Use(CORS(&CORSConfig{AllowOrigins: []string{testOriginApp}}))
	r.GET(testPath, func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func setupVaryEngineSingle(headerVal string) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set(hdrVary, headerVal)
		c.Next()
	})
	r.Use(CORS(&CORSConfig{AllowOrigins: []string{testOriginApp}}))
	r.GET(testPath, func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func TestCORS_Vary_PriorHeaders(t *testing.T) {
	r := setupVaryEngineMultiple()
	w := runVaryRequest(t, r)
	varyVals := w.Header().Values(hdrVary)
	assert.Contains(t, varyVals, "Accept-Encoding")
	assert.Contains(t, varyVals, "X-Origin-Hint")
	assert.Contains(t, varyVals, headerValOrg)
}

func TestCORS_Vary_ExistingOriginToken(t *testing.T) {
	r := setupVaryEngineSingle("Accept-Encoding, origin")
	w := runVaryRequest(t, r)
	varyVals := w.Header().Values(hdrVary)
	assert.Len(t, varyVals, 1)
	assert.Equal(t, "Accept-Encoding, origin", varyVals[0])
}

func TestCORS_Vary_ExistingWildcard(t *testing.T) {
	r := setupVaryEngineSingle("*")
	w := runVaryRequest(t, r)
	varyVals := w.Header().Values(hdrVary)
	assert.Len(t, varyVals, 1)
	assert.Equal(t, "*", varyVals[0])
}
