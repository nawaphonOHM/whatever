package server

import (
	"net/http"
)

// GetHeader returns a request header value.
func (c *Context) GetHeader(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.GetHeader(key)
	}
	return ""
}

// Header sets a response header.
func (c *Context) Header(key, value string) {
	if c.ginCtx != nil {
		c.ginCtx.Header(key, value)
	}
}

// SetHeader is an alias for Header.
func (c *Context) SetHeader(key, value string) {
	c.Header(key, value)
}

// Cookie returns a named request cookie.
func (c *Context) Cookie(name string) (string, error) {
	if c.ginCtx != nil {
		return c.ginCtx.Cookie(name)
	}
	return "", http.ErrNoCookie
}

// SetCookie adds a Set-Cookie header from an http.Cookie.
func (c *Context) SetCookie(cookie *http.Cookie) {
	if c.ginCtx == nil || cookie == nil {
		return
	}
	http.SetCookie(c.ginCtx.Writer, cookie)
}

// Status sets the HTTP response status code.
func (c *Context) Status(code int) {
	if c.ginCtx != nil {
		c.ginCtx.Status(code)
	}
}

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
