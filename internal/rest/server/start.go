package server

import (
	"context"
)

// StartREST initializes and runs the HTTP server.
// It loads OHM9969_ configuration, mounts reserved health
// endpoints, validates routes, and shuts down on signals.
func StartREST(registrations []*RestAPIRegistration) error {
	srv, err := NewFromRegistrations(registrations)
	if err != nil {
		return err
	}
	return srv.Start(context.Background())
}
