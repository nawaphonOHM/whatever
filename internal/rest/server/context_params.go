package server

// Param returns the value of the URL parameter.
func (c *Context) Param(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.Param(key)
	}
	return ""
}

// Query returns the value of the URL query parameter.
func (c *Context) Query(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.Query(key)
	}
	return ""
}

// DefaultQuery returns the value of the URL query parameter or defaultValue if
// empty.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if c.ginCtx != nil {
		return c.ginCtx.DefaultQuery(key, defaultValue)
	}
	return defaultValue
}

// QueryArray returns a slice of values for a given query key.
func (c *Context) QueryArray(key string) []string {
	if c.ginCtx != nil {
		return c.ginCtx.QueryArray(key)
	}
	return nil
}

// QueryMap returns a map of values for a given query key.
func (c *Context) QueryMap(key string) map[string]string {
	if c.ginCtx != nil {
		return c.ginCtx.QueryMap(key)
	}
	return nil
}
