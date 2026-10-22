package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContext_NilReceiver(t *testing.T) {
	var nilCtx *Context
	assert.NotNil(t, nilCtx.Context())
	assert.Equal(t, context.Background(), nilCtx.Context())
	assert.NotPanics(t, func() { nilCtx.SetRequest(nil) })
	cGinNil := &Context{ginCtx: &gin.Context{}}
	assert.NotNil(t, cGinNil.Context())
	assert.Equal(t, context.Background(), cGinNil.Context())
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
