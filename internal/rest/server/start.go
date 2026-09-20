package server

import (
	"context"

	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// StartREST initializes and runs the HTTP server.
// It loads OHM9996_ configuration, mounts reserved health
// endpoints, validates routes, and shuts down on signals.
func StartREST(bluePrint *contracts.BluePrint) error {
	srv, err := NewFromBluePrint(bluePrint)
	if err != nil {
		return err
	}
	return srv.Start(context.Background())
}
