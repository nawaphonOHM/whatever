// Package server provides internal HTTP server management, routing,
// and context handling for REST APIs.
package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// serverContextKey stores the framework Context inside gin.Context.
const serverContextKey = "_server_context"

// Context encapsulates the HTTP request lifecycle for handlers.
// It prevents direct manipulation of low-level Gin internals.
type Context struct {
	ginCtx  *gin.Context
	Request *http.Request
}

// existingContext returns a previously stored Context, if present.
func existingContext(ginCtx *gin.Context) *Context {
	val, exists := ginCtx.Get(serverContextKey)
	if !exists {
		return nil
	}
	ctx, ok := val.(*Context)
	if !ok {
		return nil
	}
	ctx.Request = ginCtx.Request
	return ctx
}

// NewContext returns a Context wrapping the given gin.Context.
func NewContext(ginCtx *gin.Context) *Context {
	if ginCtx == nil {
		return &Context{}
	}
	if ctx := existingContext(ginCtx); ctx != nil {
		return ctx
	}
	ctx := &Context{
		ginCtx:  ginCtx,
		Request: ginCtx.Request,
	}
	ginCtx.Set(serverContextKey, ctx)
	return ctx
}

// Context returns the standard Go context.Context from the underlying HTTP
// request.
func (c *Context) Context() context.Context {
	if c.ginCtx != nil && c.ginCtx.Request != nil {
		return c.ginCtx.Request.Context()
	}
	return context.Background()
}

// SetRequest updates the HTTP request on both the Context and the underlying
// gin.Context.
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

// ContentType returns the Content-Type header of the request.
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
