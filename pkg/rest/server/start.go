package server

import (
	intserver "github.com/nawaphonOHM/whatever/internal/rest/server"
)

// StartREST initializes and runs the HTTP server.
// It loads OHM9969_ configuration, mounts reserved health
// endpoints, validates routes, and shuts down on signals.
func StartREST(registrations []*RestAPIRegistration) error {
	return intserver.StartREST(registrations)
}
