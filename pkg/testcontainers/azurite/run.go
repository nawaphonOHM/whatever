package azurite

import (
	"context"

	internalazurite "github.com/nawaphonOHM/whatever/v3/internal/testcontainers/azurite"
)

// Run creates and starts an Azurite container with the provided functional options.
func Run(ctx context.Context, opts ...Option) (*Container, error) {
	return internalazurite.Run(ctx, opts...)
}
