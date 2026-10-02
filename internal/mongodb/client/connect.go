package client

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	errPingFormat          = "failed to ping mongodb: %w"
	errCreateClientFormat  = "failed to create mongodb client: %w"
	errInvalidConfigFormat = "invalid mongodb config: %w"
)

var pingClient = defaultPingClient

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, nil)
}

// verifyPingAndDisconnect verifies ping and disconnects if ping fails.
func verifyPingAndDisconnect(
	ctx context.Context,
	rawClient *mongo.Client,
) error {
	if err := pingClient(ctx, rawClient); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(
				fmt.Errorf(errPingFormat, err),
				disconnectErr,
			)
		}
		return fmt.Errorf(errPingFormat, err)
	}
	return nil
}

// resolveDatabaseName extracts the target database name from config.
func resolveDatabaseName(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	return cfg.Database
}

// verifyPingIfEnabled verifies ping connectivity if pinging is enabled in configuration.
func verifyPingIfEnabled(
	ctx context.Context,
	cfg *config.Config,
	rawClient *mongo.Client,
) error {
	if cfg == nil || !cfg.EnablePing {
		return nil
	}
	return verifyPingAndDisconnect(ctx, rawClient)
}

// attemptConnection establishes a driver client with specified TLS and optionally verifies ping.
func attemptConnection(
	ctx context.Context,
	cfg *config.Config,
	enableTLS bool,
	opts ...Option,
) (*Client, error) {
	optionsContainer := NewOptions(opts...)
	clientOptions := BuildClientOptionsWithTLS(cfg, enableTLS, optionsContainer.DriverOptions...)

	rawClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf(errCreateClientFormat, err)
	}

	if err := verifyPingIfEnabled(ctx, cfg, rawClient); err != nil {
		return nil, err
	}
	return NewClient(rawClient, resolveDatabaseName(cfg)), nil
}

// initAndPingClient performs two-phase connection: attempts unencrypted connection first,
// and falls back to TLS if the server requires TLS encryption.
func initAndPingClient(
	ctx context.Context,
	cfg *config.Config,
	opts ...Option,
) (*Client, error) {
	client, err := attemptConnection(ctx, cfg, false, opts...)
	if err == nil {
		return client, nil
	}

	if isTLSError(err) {
		return fallbackTLSAttempt(ctx, cfg, opts...)
	}

	return nil, handleConnectionError(err)
}

// Connect loads MongoDB configuration from environment variables
// (OHM9996_MONGODB_*), establishes a connection, optionally verifies ping connectivity
// if enabled, and returns a managed Client.
func Connect(ctx context.Context, opts ...Option) (*Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return ConnectWithConfig(ctx, cfg, opts...)
}

// ConnectWithConfig connects to MongoDB using the provided Config and Options.
// It verifies connectivity via Ping using the provided context if EnablePing is true.
func ConnectWithConfig(
	ctx context.Context,
	cfg *config.Config,
	opts ...Option,
) (*Client, error) {
	if cfg == nil {
		return nil, config.ErrNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf(errInvalidConfigFormat, err)
	}
	return initAndPingClient(ctx, cfg, opts...)
}
