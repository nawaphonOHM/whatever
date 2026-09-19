package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContext_AbortWithResponse(t *testing.T) {
	r := gin.New()
	r.GET("/abort", func(gc *gin.Context) {
		NewContext(gc).AbortWithResponse(BadRequest("INVALID", "bad input"))
	}, func(*gin.Context) { t.Fatal("handler ran after abort") })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/abort", nil))
	assert.Equal(t, testBadStatus, w.Code)
	assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestContext_AbortWithProblem(t *testing.T) {
	r := gin.New()
	r.GET("/problem", func(gc *gin.Context) {
		prob := NewProblemDetails(testUnauthorized, "UNAUTH", "denied")
		NewContext(gc).AbortWithProblem(prob)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/problem", nil))
	assert.Equal(t, testUnauthorized, w.Code)
}

func TestContext_AbortStatusAndNext(t *testing.T) {
	r := gin.New()
	called := false
	r.Use(func(gc *gin.Context) {
		c := NewContext(gc)
		assert.False(t, c.IsAborted())
		c.Next()
		called = true
	})
	r.GET("/status", func(gc *gin.Context) {
		c := NewContext(gc)
		c.AbortWithStatus(testBadStatus)
		assert.True(t, c.IsAborted())
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/status", nil))
	assert.Equal(t, testBadStatus, w.Code)
	assert.True(t, called)
}
