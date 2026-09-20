package server

import (
	"context"

	"github.com/nawaphonOHM/whatever/pkg/rest"
)

// StartREST initializes and runs the HTTP server.
// It loads OHM9996_ configuration, mounts reserved health
// endpoints, validates routes, and shuts down on signals.
func StartREST(registrations []*rest.RRestAPIRegistration) error {
	srv, err := NewFromRegistrations(registrations)
	if err != nil {
		return err
	}
	return srv.Start(context.Background())
}
