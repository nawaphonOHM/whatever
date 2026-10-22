package rest

import (
	"github.com/nawaphonOHM/whatever/v3/internal/rest/server"
)

// Sentinel validation errors returned by StartREST.
var (
	// ErrNilBluePrint indicates the blueprint passed to StartREST is nil.
	ErrNilBluePrint = server.ErrNilBluePrint
	// ErrReservedPath indicates a route conflicts with reserved endpoints.
	ErrReservedPath = server.ErrReservedPath
	// ErrDuplicateRoute indicates duplicate route registration.
	ErrDuplicateRoute = server.ErrDuplicateRoute
	// ErrNilRegistration indicates a nil registration entry was provided.
	ErrNilRegistration = server.ErrNilRegistration
	// ErrNilAPI indicates a nil exportable API entry was provided.
	ErrNilAPI = server.ErrNilAPI
	// ErrNilHandler indicates a nil handler was provided in an API entry.
	ErrNilHandler = server.ErrNilHandler
	// ErrInvalidMethod indicates an invalid HTTP method was provided.
	ErrInvalidMethod = server.ErrInvalidMethod
	// ErrMixedAPIGroups indicates blueprints contain mixed standard and raw API registration groups.
	ErrMixedAPIGroups = server.ErrMixedAPIGroups
	// ErrAPIGroupMismatch indicates registration group does not match configuration.
	ErrAPIGroupMismatch = server.ErrAPIGroupMismatch
	// ErrInvalidAPIGroup indicates server configuration has an invalid API group.
	ErrInvalidAPIGroup = server.ErrInvalidAPIGroup
	// ErrInvalidRegistrationGroup indicates an invalid API registration group was specified.
	ErrInvalidRegistrationGroup = server.ErrInvalidRegistrationGroup
)

// StartREST initializes and runs the HTTP server.
// It loads OHM9996_ configuration, mounts reserved health
// endpoints, validates routes, and shuts down on signals.
func StartREST(bluePrint *BluePrint) error {
	return server.StartREST(bluePrint)
}
