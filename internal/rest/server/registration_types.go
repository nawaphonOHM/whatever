package server

import (
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Handler defines route handlers returning a standardized Response.
type Handler func(*Context) response.Response

// Middleware defines the function signature for route middlewares.
type Middleware func(*Context)

// APIVersioning represents the API major version (0 = unversioned).
type APIVersioning uint
