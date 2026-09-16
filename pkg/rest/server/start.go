package server

import (
	intserver "github.com/example/go-boilerplate/internal/server"
)

// StartREST initializes and runs the HTTP server with the provided REST API registrations.
// It loads configuration with OHM9969_ prefix, mounts reserved /health and /ready endpoints,
// validates routes against collisions and duplicates, and manages graceful shutdown on SIGINT/SIGTERM.
func StartREST(registrations []*RestApiRegistration) error {
	return intserver.StartREST(registrations)
}
