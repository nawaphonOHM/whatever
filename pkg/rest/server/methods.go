package server

import (
	server2 "github.com/nawaphonOHM/whatever/internal/rest/server"
)

// HTTPMethod is the public HTTP method alias.
type HTTPMethod = server2.HTTPMethod

// Supported HTTP method constants.
const (
	GET     = server2.GET
	HEAD    = server2.HEAD
	POST    = server2.POST
	PUT     = server2.PUT
	PATCH   = server2.PATCH
	DELETE  = server2.DELETE
	OPTIONS = server2.OPTIONS
	CONNECT = server2.CONNECT
	TRACE   = server2.TRACE
)

// APIVersioning is the public API version alias.
type APIVersioning = server2.APIVersioning

// Pathz is the public path segment alias.
type Pathz = server2.Pathz
