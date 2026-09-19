package rest

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
