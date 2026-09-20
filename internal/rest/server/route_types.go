package server

import "github.com/nawaphonOHM/whatever/internal/rest/contracts"

// validatedRoute is a fully resolved route ready for engine mount.
type validatedRoute struct {
	Handler     contracts.Handler
	Path        string
	Method      string
	Middlewares []contracts.Middleware
}
