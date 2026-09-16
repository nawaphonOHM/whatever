// Package server provides the public entrypoint and declarative REST API registration types
// for bootstrapping services using the boilerplate framework.
package server

import (
	intserver "github.com/nawaphonOHM/whatever/internal/server"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Handler defines the function signature for route handlers returning a standardized Response.
type Handler = intserver.Handler

// Middleware defines the function signature for route-level middlewares.
type Middleware = intserver.Middleware

// Context encapsulates the HTTP request lifecycle for route handlers and middlewares,
// preventing accidental direct manipulation of low-level Gin engine internals.
type Context = intserver.Context

// Response represents any HTTP response that can write itself to a Gin context.
type Response = response.Response

// HTTPMethod represents supported HTTP request methods.
type HTTPMethod = intserver.HTTPMethod

const (
	GET     HTTPMethod = intserver.GET
	HEAD    HTTPMethod = intserver.HEAD
	POST    HTTPMethod = intserver.POST
	PUT     HTTPMethod = intserver.PUT
	PATCH   HTTPMethod = intserver.PATCH
	DELETE  HTTPMethod = intserver.DELETE
	OPTIONS HTTPMethod = intserver.OPTIONS
	CONNECT HTTPMethod = intserver.CONNECT
	TRACE   HTTPMethod = intserver.TRACE
)

// ApiVersioning represents the API major version (e.g. 1 for v1). 0 indicates unversioned.
type ApiVersioning = intserver.ApiVersioning

// Pathz represents a URL path segment.
type Pathz = intserver.Pathz

// ExportableApi defines a single API route endpoint with method, path, middlewares, and handler.
type ExportableApi = intserver.ExportableApi

// RestApiRegistration groups multiple exportable APIs under a common prefix and version.
type RestApiRegistration = intserver.RestApiRegistration

// Sentinel validation errors.
var (
	ErrReservedPath    = intserver.ErrReservedPath
	ErrDuplicateRoute  = intserver.ErrDuplicateRoute
	ErrNilRegistration = intserver.ErrNilRegistration
	ErrNilAPI          = intserver.ErrNilAPI
	ErrNilHandler      = intserver.ErrNilHandler
	ErrInvalidMethod   = intserver.ErrInvalidMethod
)
