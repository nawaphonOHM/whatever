package server

// ShouldBind selects a binding engine from Method and Content-Type.
func (c *Context) ShouldBind(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBind(obj)
	}
	return nil
}

// ShouldBindJSON binds a JSON request body.
func (c *Context) ShouldBindJSON(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindJSON(obj)
	}
	return nil
}

// ShouldBindQuery binds query parameters.
func (c *Context) ShouldBindQuery(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindQuery(obj)
	}
	return nil
}

// ShouldBindURI binds URI parameters.
func (c *Context) ShouldBindURI(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindUri(obj)
	}
	return nil
}

// ShouldBindHeader binds request headers.
func (c *Context) ShouldBindHeader(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.ShouldBindHeader(obj)
	}
	return nil
}

// Bind selects a binding engine from Method and Content-Type.
func (c *Context) Bind(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.Bind(obj)
	}
	return nil
}

// BindJSON binds a JSON request body.
func (c *Context) BindJSON(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindJSON(obj)
	}
	return nil
}

// BindQuery binds query parameters.
func (c *Context) BindQuery(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindQuery(obj)
	}
	return nil
}

// BindURI binds URI parameters.
func (c *Context) BindURI(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindUri(obj)
	}
	return nil
}

// BindHeader binds request headers.
func (c *Context) BindHeader(obj any) error {
	if c.ginCtx != nil {
		return c.ginCtx.BindHeader(obj)
	}
	return nil
}

// PostForm returns a form value from a POST, PATCH, or PUT body.
func (c *Context) PostForm(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.PostForm(key)
	}
	return ""
}

// DefaultPostForm returns a form value or defaultValue if it is empty.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if c.ginCtx != nil {
		return c.ginCtx.DefaultPostForm(key, defaultValue)
	}
	return defaultValue
}

// PostFormArray returns values for a form key.
func (c *Context) PostFormArray(key string) []string {
	if c.ginCtx != nil {
		return c.ginCtx.PostFormArray(key)
	}
	return nil
}

// PostFormMap returns values for a form key as a map.
func (c *Context) PostFormMap(key string) map[string]string {
	if c.ginCtx != nil {
		return c.ginCtx.PostFormMap(key)
	}
	return nil
}
