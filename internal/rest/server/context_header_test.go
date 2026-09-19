package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContext_HeadersAndCookies(t *testing.T) {
	r := gin.New()
	r.GET("/headers", func(gc *gin.Context) {
		c := newContext(gc)
		assert.Equal(t, "Bearer xyz", c.GetHeader("Authorization"))
		cookie, err := c.Cookie("session")
		assert.NoError(t, err)
		assert.Equal(t, "session-val", cookie)
		c.Header("X-Custom", testValue)
		c.SetHeader("X-Alias", testValue)
		c.SetCookie(&http.Cookie{Name: "remember", Value: testValue})
		c.Status(testAccepted)
	})
	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("Authorization", "Bearer xyz")
	req.AddCookie(&http.Cookie{Name: "session", Value: "session-val"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, testAccepted, w.Code)
	assert.Equal(t, testValue, w.Header().Get("X-Custom"))
	assert.Contains(t, w.Header().Get("Set-Cookie"), "remember=value")
}
