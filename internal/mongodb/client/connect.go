package client

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

// verifyPingAndDisconnect verifies ping and disconnects if ping fails.
func verifyPingAndDisconnect(
	ctx context.Context,
	rawClient *mongo.Client,
) error {
	if err := rawClient.Ping(ctx, nil); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(
				fmt.Errorf("failed to ping mongodb: %w", err),
				disconnectErr,
			)
		}
		return fmt.Errorf("failed to ping mongodb: %w", err)
	}
	return nil
}

// initAndPingClient instantiates the driver client and verifies ping
// connectivity.
func initAndPingClient(
	ctx context.Context,
	cfg *config.Config,
	opts ...Option,
) (*Client, error) {
	optionsContainer := NewOptions(opts...)
	clientOptions := BuildClientOptions(cfg, optionsContainer.DriverOptions...)

	rawClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	if err := verifyPingAndDisconnect(ctx, rawClient); err != nil {
		return nil, err
	}
	return NewClient(rawClient, cfg.Database), nil
}

// Connect loads MongoDB configuration from environment variables
// (OHM9996_MONGODB_*), establishes a connection, verifies ping connectivity,
// and returns a managed Client.
func Connect(ctx context.Context, opts ...Option) (*Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return ConnectWithConfig(ctx, cfg, opts...)
}

// ConnectWithConfig connects to MongoDB using the provided Config and Options.
// It verifies connectivity via Ping using the provided context.
func ConnectWithConfig(
	ctx context.Context,
	cfg *config.Config,
	opts ...Option,
) (*Client, error) {
	if cfg == nil {
		return nil, config.ErrNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid mongodb config: %w", err)
	}
	return initAndPingClient(ctx, cfg, opts...)
}
