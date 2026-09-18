package server

// ShouldBind checks the Method and Content-Type to select a binding engine
// automatically.
func (c *Context) ShouldBind(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBind(obj)
	}
	return nil
}

// ShouldBindJSON is a shortcut for c.ShouldBindWith(obj, binding.JSON).
func (c *Context) ShouldBindJSON(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindJSON(obj)
	}
	return nil
}

// ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
func (c *Context) ShouldBindQuery(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindQuery(obj)
	}
	return nil
}

// ShouldBindURI is a shortcut for c.ShouldBindWith(obj, binding.Uri).
func (c *Context) ShouldBindURI(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindUri(obj)
	}
	return nil
}

// ShouldBindHeader is a shortcut for c.ShouldBindWith(obj, binding.Header).
func (c *Context) ShouldBindHeader(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindHeader(obj)
	}
	return nil
}

// Bind checks the Method and Content-Type to select a binding engine
// automatically.
func (c *Context) Bind(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.Bind(obj)
	}
	return nil
}

// BindJSON is a shortcut for c.BindWith(obj, binding.JSON).
func (c *Context) BindJSON(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindJSON(obj)
	}
	return nil
}

// BindQuery is a shortcut for c.BindWith(obj, binding.Query).
func (c *Context) BindQuery(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindQuery(obj)
	}
	return nil
}

// BindURI is a shortcut for c.BindWith(obj, binding.Uri).
func (c *Context) BindURI(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindUri(obj)
	}
	return nil
}

// BindHeader is a shortcut for c.BindWith(obj, binding.Header).
func (c *Context) BindHeader(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindHeader(obj)
	}
	return nil
}
