package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContext_LifecycleAndReentrancy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	gc, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	gc.Request = req
	c1 := newContext(gc)
	c2 := newContext(gc)
	assert.Same(t, c1, c2)
	assert.Same(t, gc, c1.ginCtx)
	assert.Same(t, req, c1.Request)
}

func TestContext_NilContext(t *testing.T) {
	c := newContext(nil)
	assert.NotNil(t, c)
}

func TestContext_MetadataAndRequest(t *testing.T) {
	r := gin.New()
	var got *http.Request
	r.POST("/meta/:id", func(gc *gin.Context) {
		c := newContext(gc)
		assert.NotEmpty(t, c.ClientIP())
		assert.Equal(t, "application/json", c.ContentType())
		assert.Equal(t, "/meta/:id", c.FullPath())
		got = gc.Request.WithContext(context.Background())
		c.SetRequest(got)
		c.Status(testStatus)
	})
	req := httptest.NewRequest(http.MethodPost, "/meta/123", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, testStatus, w.Code)
	assert.NotNil(t, got)
}

func TestContext_NilSafety(t *testing.T) {
	c := &Context{}
	assert.Empty(t, c.ClientIP())
	assert.Empty(t, c.ContentType())
	assert.Empty(t, c.FullPath())
	assert.False(t, c.IsAborted())
	assert.NoError(t, c.ShouldBind(nil))
	assert.NoError(t, c.Bind(nil))
	assert.Empty(t, c.PostForm(testKey))
	assert.Equal(t, testDefault, c.DefaultPostForm(testKey, testDefault))
	assert.Equal(t, context.Background(), c.Context())
}
