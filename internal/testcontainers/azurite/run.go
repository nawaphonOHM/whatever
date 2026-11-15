package azurite

import (
	"context"
	"errors"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	tcazurite "github.com/testcontainers/testcontainers-go/modules/azure/azurite"
)

// Run creates and starts an Azurite container with the provided functional options.
func Run(ctx context.Context, opts ...Option) (*Container, error) {
	return runWithRunner(ctx, tcazurite.Run, opts...)
}

func runWithRunner(
	ctx context.Context,
	runner func(context.Context, string, ...testcontainers.ContainerCustomizer) (*tcazurite.Container, error),
	opts ...Option,
) (*Container, error) {
	options := NewOptions(opts...)
	raw, err := runner(ctx, options.Image, buildCustomizers(options)...)
	if err != nil {
		return nil, fmt.Errorf("failed to start azurite container: %w", cleanupFailedRun(ctx, raw, err))
	}
	if raw == nil {
		return nil, errors.New("azurite startup returned no container")
	}
	return NewContainer(raw), nil
}

func cleanupFailedRun(ctx context.Context, raw *tcazurite.Container, startErr error) error {
	if raw == nil {
		return startErr
	}
	cleanupErr := raw.Terminate(ctx)
	if cleanupErr == nil {
		return startErr
	}
	return errors.Join(startErr, fmt.Errorf("failed to terminate partially started azurite container: %w", cleanupErr))
}
