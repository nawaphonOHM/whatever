package rest

import (
	"time"
)

// Context defines the read-only request execution context supplied to
// handlers and middlewares. It provides getter and binding access to the
// incoming HTTP request without exposing internal mutating operations.
type Context interface {
	// ClientIP returns the client's IP address.
	ClientIP() string
	// ContentType returns the request Content-Type header.
	ContentType() string
	// FullPath returns the matched route full path.
	FullPath() string

	// Param returns the value of a URL parameter.
	Param(string) string
	// Query returns the value of a URL query parameter.
	Query(string) string
	// DefaultQuery returns a query value or defaultValue if empty.
	DefaultQuery(string, string) string
	// QueryArray returns values for a query key as a slice.
	QueryArray(string) []string
	// QueryMap returns values for a query key as a map.
	QueryMap(string) map[string]string

	// PostForm returns a form value from a POST, PATCH, or PUT body.
	PostForm(string) string
	// DefaultPostForm returns a form value or defaultValue if empty.
	DefaultPostForm(string, string) string
	// PostFormArray returns values for a form key as a slice.
	PostFormArray(string) []string
	// PostFormMap returns values for a form key as a map.
	PostFormMap(string) map[string]string

	// GetHeader returns a request header value.
	GetHeader(string) string
	// Cookie returns a named request cookie.
	Cookie(string) (string, error)

	// ShouldBind selects a binding engine from Method and Content-Type.
	ShouldBind(any) error
	// ShouldBindJSON binds a JSON request body.
	ShouldBindJSON(any) error
	// ShouldBindQuery binds query parameters.
	ShouldBindQuery(any) error
	// ShouldBindURI binds URI parameters.
	ShouldBindURI(any) error
	// ShouldBindHeader binds request headers.
	ShouldBindHeader(any) error
	// Bind selects a binding engine from Method and Content-Type.
	Bind(any) error
	// BindJSON binds a JSON request body.
	BindJSON(any) error
	// BindQuery binds query parameters.
	BindQuery(any) error
	// BindURI binds URI parameters.
	BindURI(any) error
	// BindHeader binds request headers.
	BindHeader(any) error

	// Get returns the value for a key and whether it exists.
	Get(string) (any, bool)
	// MustGet returns the value for a key or panics if missing.
	MustGet(string) any
	// GetString returns a context value as a string.
	GetString(string) string
	// GetBool returns a context value as a boolean.
	GetBool(string) bool
	// GetInt returns a context value as an integer.
	GetInt(string) int
	// GetInt64 returns a context value as int64.
	GetInt64(string) int64
	// GetUint returns a context value as uint.
	GetUint(string) uint
	// GetUint64 returns a context value as uint64.
	GetUint64(string) uint64
	// GetFloat64 returns a context value as float64.
	GetFloat64(string) float64
	// GetTime returns a context value as time.Time.
	GetTime(string) time.Time
	// GetDuration returns a context value as time.Duration.
	GetDuration(string) time.Duration
	// GetStringSlice returns a context value as a string slice.
	GetStringSlice(string) []string
	// GetStringMap returns a context value as a map of interfaces.
	GetStringMap(string) map[string]any
	// GetStringMapString returns a context value as a string map.
	GetStringMapString(string) map[string]string
	// GetStringMapStringSlice returns a context value as map of string slices.
	GetStringMapStringSlice(string) map[string][]string
}
