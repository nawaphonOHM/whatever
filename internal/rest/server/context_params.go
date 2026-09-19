package server

// Param returns the value of a URL parameter.
func (c *Context) Param(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.Param(key)
	}
	return ""
}

// Query returns the value of a URL query parameter.
func (c *Context) Query(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.Query(key)
	}
	return ""
}

// DefaultQuery returns a query value or defaultValue if it is empty.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if c.ginCtx != nil {
		return c.ginCtx.DefaultQuery(key, defaultValue)
	}
	return defaultValue
}

// QueryArray returns values for a query key.
func (c *Context) QueryArray(key string) []string {
	if c.ginCtx != nil {
		return c.ginCtx.QueryArray(key)
	}
	return nil
}

// QueryMap returns values for a query key as a map.
func (c *Context) QueryMap(key string) map[string]string {
	if c.ginCtx != nil {
		return c.ginCtx.QueryMap(key)
	}
	return nil
}
