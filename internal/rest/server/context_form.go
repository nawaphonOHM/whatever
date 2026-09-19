package server

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
