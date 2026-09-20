package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

const serverContextKey = "_server_context"

// Context implements contracts.Context wrapping gin.Context.
type Context struct {
	ginCtx  *gin.Context
	Request *http.Request
}

// Compile-time assertion that Context implements contracts.Context.
var _ contracts.Context = (*Context)(nil)

func existingContext(ginCtx *gin.Context) *Context {
	val, exists := ginCtx.Get(serverContextKey)
	if ctx, ok := val.(*Context); exists && ok {
		ctx.Request = ginCtx.Request
		return ctx
	}
	return nil
}

// newContext returns a concrete Context wrapping the given gin.Context.
func newContext(ginCtx *gin.Context) *Context {
	if ginCtx == nil {
		return &Context{}
	}
	if ctx := existingContext(ginCtx); ctx != nil {
		return ctx
	}
	ctx := &Context{ginCtx: ginCtx, Request: ginCtx.Request}
	ginCtx.Set(serverContextKey, ctx)
	return ctx
}

// SetRequest updates the request on both Context and the Gin context.
func (c *Context) SetRequest(r *http.Request) {
	c.Request = r
	if c.ginCtx != nil {
		c.ginCtx.Request = r
	}
}

// ClientIP returns the client's IP address.
func (c *Context) ClientIP() string {
	if c.ginCtx != nil {
		return c.ginCtx.ClientIP()
	}
	return ""
}

// ContentType returns the request Content-Type header.
func (c *Context) ContentType() string {
	if c.ginCtx != nil {
		return c.ginCtx.ContentType()
	}
	return ""
}

// FullPath returns the matched route full path.
func (c *Context) FullPath() string {
	if c.ginCtx != nil {
		return c.ginCtx.FullPath()
	}
	return ""
}
