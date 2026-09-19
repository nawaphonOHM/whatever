// Package rest provides standardized API registration contracts and
// JSON response envelopes according to RFC 9457 Problem Details.
//
//revive:disable:max-public-structs // required by the unified public API
package rest

// APIVersioning represents the API major version (0 = unversioned).
type APIVersioning uint

// Pathz represents a URL path segment.
type Pathz string

// Handler defines route handlers returning a standardized Response.
type Handler func(Context) Response

// Middleware defines the function signature for route middlewares.
type Middleware func(Context)

// ExportableAPI defines a single API route endpoint.
type ExportableAPI struct {
	Path       Pathz
	Handler    Handler
	Middleware []Middleware
	Method     HTTPMethod
}

// RRestAPIRegistration groups APIs under a common prefix and version.
type RRestAPIRegistration struct {
	Prefix  Pathz
	Apis    []*ExportableAPI
	Version APIVersioning
}
