package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_NilSafety(t *testing.T) {
	c := &Context{}
	assert.Empty(t, c.ClientIP())
	assert.Empty(t, c.ContentType())
	assert.Empty(t, c.FullPath())
	assert.NotNil(t, c.Context())
	assert.False(t, c.IsAborted())
	assert.NoError(t, c.ShouldBind(nil))
	assert.NoError(t, c.Bind(nil))
	assert.Empty(t, c.PostForm(testKey))
	assert.Equal(t, testDefault, c.DefaultPostForm(testKey, testDefault))
}

func TestContext_NilQueriesAndHeaders(t *testing.T) {
	c := &Context{}
	assert.Empty(t, c.Query(testKey))
	assert.Equal(t, testDefault, c.DefaultQuery(testKey, testDefault))
	assert.Empty(t, c.GetHeader(testKey))
	_, err := c.Cookie(testKey)
	assert.ErrorIs(t, err, http.ErrNoCookie)
}

func TestContext_NilState(t *testing.T) {
	c := &Context{}
	c.Set(testKey, testValue)
	value, exists := c.Get(testKey)
	assert.Nil(t, value)
	assert.False(t, exists)
	assert.Panics(t, func() { c.MustGet(testKey) })
}

func TestContext_NilFlow(t *testing.T) {
	c := &Context{}
	assert.NotNil(t, c)
	c.SetHeader(testKey, testValue)
	c.SetCookie(nil)
	c.Status(testStatus)
	c.Next()
	c.Abort()
	c.AbortWithStatus(testBadStatus)
	c.AbortWithResponse(nil)
	c.AbortWithProblem(nil)
}
