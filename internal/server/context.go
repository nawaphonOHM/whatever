package server

import (
	"context"
	"net/http"
	"time"

	"github.com/example/go-boilerplate/pkg/rest/response"
	"github.com/gin-gonic/gin"
)

// Context encapsulates the HTTP request lifecycle for route handlers and middlewares,
// preventing accidental direct manipulation of low-level Gin engine internals.
type Context struct {
	ginCtx  *gin.Context
	Request *http.Request
}

// NewContext returns a Context wrapping the given gin.Context.
func NewContext(ginCtx *gin.Context) *Context {
	if ginCtx == nil {
		return &Context{}
	}
	if c, exists := ginCtx.Get("_server_context"); exists {
		if ctx, ok := c.(*Context); ok {
			ctx.Request = ginCtx.Request
			return ctx
		}
	}
	ctx := &Context{
		ginCtx:  ginCtx,
		Request: ginCtx.Request,
	}
	ginCtx.Set("_server_context", ctx)
	return ctx
}

// Context returns the standard Go context.Context from the underlying HTTP request.
func (c *Context) Context() context.Context {
	if c.ginCtx != nil && c.ginCtx.Request != nil {
		return c.ginCtx.Request.Context()
	}
	return context.Background()
}

// SetRequest updates the HTTP request on both the Context and the underlying gin.Context.
func (c *Context) SetRequest(r *http.Request) {
	c.Request = r
	if c.ginCtx != nil {
		c.ginCtx.Request = r
	}
}

// ClientIP returns the client's IP address.
func (c *Context) ClientIP() string {
	if c.ginCtx != nil {
		return c.ginCtx.ClientIP()
	}
	return ""
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	if c.ginCtx != nil {
		return c.ginCtx.ContentType()
	}
	return ""
}

// FullPath returns the matched route full path.
func (c *Context) FullPath() string {
	if c.ginCtx != nil {
		return c.ginCtx.FullPath()
	}
	return ""
}

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

// DefaultQuery returns the value of the URL query parameter or defaultValue if empty.
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

// PostForm returns the specified key from a POST, PATCH, or PUT body.
func (c *Context) PostForm(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.PostForm(key)
	}
	return ""
}

// DefaultPostForm returns the specified key from a POST body or defaultValue if empty.
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

// GetHeader returns the value of the specified request header.
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

// SetHeader is an alias for Header to set a response header.
func (c *Context) SetHeader(key, value string) {
	c.Header(key, value)
}

// Cookie returns the named cookie provided in the request.
func (c *Context) Cookie(name string) (string, error) {
	if c.ginCtx != nil {
		return c.ginCtx.Cookie(name)
	}
	return "", http.ErrNoCookie
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if c.ginCtx != nil {
		c.ginCtx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	}
}

// Status sets the HTTP response status code.
func (c *Context) Status(code int) {
	if c.ginCtx != nil {
		c.ginCtx.Status(code)
	}
}

// ShouldBind checks the Method and Content-Type to select a binding engine automatically.
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

// ShouldBindUri is a shortcut for c.ShouldBindWith(obj, binding.Uri).
func (c *Context) ShouldBindUri(obj any) error {
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

// Bind checks the Method and Content-Type to select a binding engine automatically.
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

// BindUri is a shortcut for c.BindWith(obj, binding.Uri).
func (c *Context) BindUri(obj any) error {
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

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) any {
	if c.ginCtx != nil {
		return c.ginCtx.MustGet(key)
	}
	panic("context is nil")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) string {
	if c.ginCtx != nil {
		return c.ginCtx.GetString(key)
	}
	return ""
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) bool {
	if c.ginCtx != nil {
		return c.ginCtx.GetBool(key)
	}
	return false
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) int {
	if c.ginCtx != nil {
		return c.ginCtx.GetInt(key)
	}
	return 0
}

// GetInt64 returns the value associated with the key as an int64.
func (c *Context) GetInt64(key string) int64 {
	if c.ginCtx != nil {
		return c.ginCtx.GetInt64(key)
	}
	return 0
}

// GetUint returns the value associated with the key as a uint.
func (c *Context) GetUint(key string) uint {
	if c.ginCtx != nil {
		return c.ginCtx.GetUint(key)
	}
	return 0
}

// GetUint64 returns the value associated with the key as a uint64.
func (c *Context) GetUint64(key string) uint64 {
	if c.ginCtx != nil {
		return c.ginCtx.GetUint64(key)
	}
	return 0
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) float64 {
	if c.ginCtx != nil {
		return c.ginCtx.GetFloat64(key)
	}
	return 0
}

// GetTime returns the value associated with the key as time.Time.
func (c *Context) GetTime(key string) time.Time {
	if c.ginCtx != nil {
		return c.ginCtx.GetTime(key)
	}
	return time.Time{}
}

// GetDuration returns the value associated with the key as time.Duration.
func (c *Context) GetDuration(key string) time.Duration {
	if c.ginCtx != nil {
		return c.ginCtx.GetDuration(key)
	}
	return 0
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) []string {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringSlice(key)
	}
	return nil
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) map[string]any {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringMap(key)
	}
	return nil
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) map[string]string {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringMapString(key)
	}
	return nil
}

// GetStringMapStringSlice returns the value associated with the key as a map of string slices.
func (c *Context) GetStringMapStringSlice(key string) map[string][]string {
	if c.ginCtx != nil {
		return c.ginCtx.GetStringMapStringSlice(key)
	}
	return nil
}

// Next executes the pending handlers in the chain inside the calling middleware.
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

// AbortWithStatus calls Abort and writes the headers with the specified status code.
func (c *Context) AbortWithStatus(code int) {
	if c.ginCtx != nil {
		c.ginCtx.AbortWithStatus(code)
	}
}

// AbortWithResponse aborts the handler chain and serializes the provided Response.
func (c *Context) AbortWithResponse(resp response.Response) {
	if c.ginCtx != nil {
		c.ginCtx.Abort()
		if resp != nil {
			resp.Write(c.ginCtx)
		}
	}
}

// AbortWithProblem aborts the handler chain and serializes the provided RFC 9457 ProblemDetails.
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
