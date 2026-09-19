package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
)

func TestContext_AbortWithResponse(t *testing.T) {
	r := gin.New()
	r.GET("/abort", func(gc *gin.Context) {
		c := newContext(gc)
		c.AbortWithResponse(rest.BadRequest("INVALID", "bad input"))
	}, func(*gin.Context) { t.Fatal("handler ran after abort") })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/abort", nil))
	assert.Equal(t, testBadStatus, w.Code)
	assert.Equal(t, rest.MediaTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestContext_AbortWithProblem(t *testing.T) {
	r := gin.New()
	r.GET("/problem", func(gc *gin.Context) {
		prob := rest.NewProblemDetails(testUnauthorized, "UNAUTH", "denied")
		newContext(gc).AbortWithProblem(prob)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/problem", nil))
	assert.Equal(t, testUnauthorized, w.Code)
}

func TestContext_AbortStatusAndNext(t *testing.T) {
	r := gin.New()
	called := false
	r.Use(func(gc *gin.Context) {
		c := newContext(gc)
		assert.False(t, c.IsAborted())
		c.Next()
		called = true
	})
	r.GET("/status", func(gc *gin.Context) {
		c := newContext(gc)
		c.AbortWithStatus(testBadStatus)
		assert.True(t, c.IsAborted())
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/status", nil))
	assert.Equal(t, testBadStatus, w.Code)
	assert.True(t, called)
}
