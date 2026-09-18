package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	customHeaderName = "X-Custom"
	customHeaderVal  = "custom-value"
	respHeaderName   = "X-Response"
	respHeaderVal    = "resp-value"
	cookieName       = "session"
	cookieVal        = "abc"
	cookieMaxAge     = 3600
	cookiePath       = "/"
	cookieDomain     = "example.com"
	headersPath      = "/test-headers"
)

// applyHeaderCookieWrites writes response header and cookie.
func applyHeaderCookieWrites(c *Context) {
	c.SetHeader(respHeaderName, respHeaderVal)
	c.SetCookie(&http.Cookie{
		Name:     cookieName,
		Value:    cookieVal,
		MaxAge:   cookieMaxAge,
		Path:     cookiePath,
		Domain:   cookieDomain,
		Secure:   true,
		HttpOnly: true,
	})
	c.Status(http.StatusOK)
}

// runHeaderCookieRequest exercises the headers test route.
func runHeaderCookieRequest(
	r *gin.Engine,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, headersPath, nil)
	req.Header.Set(customHeaderName, customHeaderVal)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestContext_HeadersAndCookies covers header and cookie helpers.
func TestContext_HeadersAndCookies(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	r := gin.New()
	var headerVal string
	r.GET(headersPath, func(gc *gin.Context) {
		c := NewContext(gc)
		headerVal = c.GetHeader(customHeaderName)
		_, err := c.Cookie(cookieName)
		assert.ErrorIs(t, err, http.ErrNoCookie)
		applyHeaderCookieWrites(c)
	})

	// Act
	w := runHeaderCookieRequest(r)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, customHeaderVal, headerVal)
	assert.Equal(t, respHeaderVal, w.Header().Get(respHeaderName))
	assert.Contains(t, w.Header().Get("Set-Cookie"), cookieName)
}
