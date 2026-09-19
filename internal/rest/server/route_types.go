package server

import "github.com/nawaphonOHM/whatever/pkg/rest"

// validatedRoute is a fully resolved route ready for engine mount.
type validatedRoute struct {
	Handler     rest.Handler
	Path        string
	Method      string
	Middlewares []rest.Middleware
}
