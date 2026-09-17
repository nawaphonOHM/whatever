// Package server provides the public entrypoint and declarative REST API registration types
// for bootstrapping services using the boilerplate framework.
package server

import (
	server2 "github.com/nawaphonOHM/whatever/internal/rest/server"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Handler defines the function signature for route handlers returning a standardized Response.
type Handler = server2.Handler

// Middleware defines the function signature for route-level middlewares.
type Middleware = server2.Middleware

// Context encapsulates the HTTP request lifecycle for route handlers and middlewares,
// preventing accidental direct manipulation of low-level Gin engine internals.
type Context = server2.Context

// Response represents any HTTP response that can write itself to a Gin context.
type Response = response.Response

// HTTPMethod represents supported HTTP request methods.
type HTTPMethod = server2.HTTPMethod

// Supported HTTP method constants.
const (
	GET     HTTPMethod = server2.GET
	HEAD    HTTPMethod = server2.HEAD
	POST    HTTPMethod = server2.POST
	PUT     HTTPMethod = server2.PUT
	PATCH   HTTPMethod = server2.PATCH
	DELETE  HTTPMethod = server2.DELETE
	OPTIONS HTTPMethod = server2.OPTIONS
	CONNECT HTTPMethod = server2.CONNECT
	TRACE   HTTPMethod = server2.TRACE
)

// APIVersioning represents the API major version (e.g. 1 for v1). 0 indicates unversioned.
type APIVersioning = server2.APIVersioning

// Pathz represents a URL path segment.
type Pathz = server2.Pathz

// ExportableAPI defines a single API route endpoint with method, path, middlewares, and handler.
type ExportableAPI = server2.ExportableAPI

// RestAPIRegistration groups multiple exportable APIs under a common prefix and version.
type RestAPIRegistration = server2.RestAPIRegistration

// Sentinel validation errors.
var (
	ErrReservedPath    = server2.ErrReservedPath
	ErrDuplicateRoute  = server2.ErrDuplicateRoute
	ErrNilRegistration = server2.ErrNilRegistration
	ErrNilAPI          = server2.ErrNilAPI
	ErrNilHandler      = server2.ErrNilHandler
	ErrInvalidMethod   = server2.ErrInvalidMethod
)
