package rest

// Set stores a key/value pair exclusively for this context.
func (c *Context) Set(key string, value any) {
	if c.ginCtx != nil {
		c.ginCtx.Set(key, value)
	}
}

// Get returns the value for a key and whether it exists.
func (c *Context) Get(key string) (value any, exists bool) {
	if c.ginCtx != nil {
		return c.ginCtx.Get(key)
	}
	return nil, false
}

// MustGet returns the value for a key or panics if it is missing.
func (c *Context) MustGet(key string) any {
	if c.ginCtx != nil {
		return c.ginCtx.MustGet(key)
	}
	panic("context is nil")
}

// Next executes pending handlers in the chain.
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

// AbortWithStatus aborts and writes the specified status code.
func (c *Context) AbortWithStatus(code int) {
	if c.ginCtx != nil {
		c.ginCtx.AbortWithStatus(code)
	}
}

// AbortWithResponse aborts and serializes the provided Response.
func (c *Context) AbortWithResponse(resp Response) {
	if c.ginCtx == nil {
		return
	}
	c.ginCtx.Abort()
	if resp != nil {
		resp.Write(c.ginCtx)
	}
}

// AbortWithProblem aborts and serializes the provided ProblemDetails.
func (c *Context) AbortWithProblem(prob *ProblemDetails) {
	c.AbortWithResponse(prob)
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	if c.ginCtx != nil {
		return c.ginCtx.IsAborted()
	}
	return false
}
