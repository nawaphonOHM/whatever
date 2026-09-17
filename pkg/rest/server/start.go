// Package server provides the public entrypoint and declarative REST API registration types
// for bootstrapping services using the boilerplate framework.
package server

import (
	intserver "github.com/nawaphonOHM/whatever/internal/server"
)

// StartREST initializes and runs the HTTP server with the provided REST API registrations.
// It loads configuration with OHM9969_ prefix, mounts reserved /health and /ready endpoints,
// validates routes against collisions and duplicates, and manages graceful shutdown on SIGINT/SIGTERM.
func StartREST(registrations []*RestAPIRegistration) error {
	return intserver.StartREST(registrations)
}
