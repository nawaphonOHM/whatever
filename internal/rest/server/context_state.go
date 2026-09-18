package server

import (
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Set stores a new key/value pair exclusively for this context.
func (c *Context) Set(key string, value any) {
	if c.ginCtx != nil {
		c.ginCtx.Set(key, value)
	}
}

// Get returns the value for the given key and whether it exists.
func (c *Context) Get(key string) (value any, exists bool) {
	if c.ginCtx != nil {
		return c.ginCtx.Get(key)
	}
	return nil, false
}

// MustGet returns the value for the given key if it exists, otherwise it
// panics.
func (c *Context) MustGet(key string) any {
	if c.ginCtx != nil {
		return c.ginCtx.MustGet(key)
	}
	panic("context is nil")
}

// Next executes the pending handlers in the chain inside the calling
// middleware.
func (c *Context) Next() {
	if c.ginCtx != nil {
		c.ginCtx.Next()
	}
}

// Abort prevents pending handlers from being called.
func (c *Context) Abort() {
	if c.ginCtx != nil {
		c.ginCtx.Abort()
	}
}

// AbortWithStatus calls Abort and writes the headers with the specified status
// code.
func (c *Context) AbortWithStatus(code int) {
	if c.ginCtx != nil {
		c.ginCtx.AbortWithStatus(code)
	}
}

// AbortWithResponse aborts the handler chain and serializes the provided
// Response.
func (c *Context) AbortWithResponse(resp response.Response) {
	if c.ginCtx == nil {
		return
	}
	c.ginCtx.Abort()
	if resp != nil {
		resp.Write(c.ginCtx)
	}
}

// AbortWithProblem aborts the handler chain and serializes the provided RFC
// 9457 ProblemDetails.
func (c *Context) AbortWithProblem(prob *response.ProblemDetails) {
	c.AbortWithResponse(prob)
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	if c.ginCtx != nil {
		return c.ginCtx.IsAborted()
	}
	return false
}
