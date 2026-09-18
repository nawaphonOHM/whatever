package server

// PostForm returns the specified key from a POST, PATCH, or PUT body.
func (c *Context) PostForm(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.PostForm(key)
	}
	return ""
}

// DefaultPostForm returns the specified key from a POST body or defaultValue if
// empty.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if c.ginCtx != nil {
		return c.ginCtx.DefaultPostForm(key, defaultValue)
	}
	return defaultValue
}

// PostFormArray returns a slice of strings for a given form key.
func (c *Context) PostFormArray(key string) []string {
	if c.ginCtx != nil {
		return c.ginCtx.PostFormArray(key)
	}
	return nil
}

// PostFormMap returns a map of strings for a given form key.
func (c *Context) PostFormMap(key string) map[string]string {
	if c.ginCtx != nil {
		return c.ginCtx.PostFormMap(key)
	}
	return nil
}
