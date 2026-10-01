package mongodb

import (
	"context"
	"fmt"

	tcmongodb "github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// Run creates and starts a MongoDB container with the provided functional options.
func Run(ctx context.Context, opts ...Option) (*Container, error) {
	options := NewOptions(opts...)
	customizers := buildCustomizers(options)
	raw, err := tcmongodb.Run(ctx, options.Image, customizers...)
	if err != nil {
		return nil, fmt.Errorf("failed to start mongodb container: %w", err)
	}
	return NewContainer(raw, options), nil
}
