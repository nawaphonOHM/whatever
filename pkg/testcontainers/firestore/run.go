package firestore

import (
	"context"

	"github.com/nawaphonOHM/whatever/v2/internal/testcontainers/firestore"
)

// Run creates and starts a Firestore container with the provided functional options.
func Run(ctx context.Context, opts ...Option) (*Container, error) {
	return firestore.Run(ctx, opts...)
}
