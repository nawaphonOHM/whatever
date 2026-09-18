// Package server provides the public entrypoint and declarative
// REST API registration types for bootstrapping services.
package server

import (
	server2 "github.com/nawaphonOHM/whatever/internal/rest/server"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Handler defines public route handlers.
type Handler func(*Context) response.Response

// Middleware defines public route middlewares.
type Middleware func(*Context)

// Context is the public request context alias.
type Context = server2.Context
