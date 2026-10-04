package mongodb

import (
	"context"

	"github.com/nawaphonOHM/whatever/v2/internal/testcontainers/mongodb"
)

// Run creates and starts a MongoDB container with the provided functional options.
func Run(ctx context.Context, opts ...Option) (*Container, error) {
	return mongodb.Run(ctx, opts...)
}
