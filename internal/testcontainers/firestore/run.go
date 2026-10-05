package firestore

import (
	"context"
	"fmt"

	tcfirestore "github.com/testcontainers/testcontainers-go/modules/gcloud/firestore"
)

// Run creates and starts a Firestore container with the provided functional options.
func Run(ctx context.Context, opts ...Option) (*Container, error) {
	options := NewOptions(opts...)
	customizers := buildCustomizers(options)
	raw, err := tcfirestore.Run(ctx, options.Image, customizers...)
	if err != nil {
		return nil, fmt.Errorf("failed to start firestore container: %w", err)
	}
	return NewContainer(raw, options), nil
}
