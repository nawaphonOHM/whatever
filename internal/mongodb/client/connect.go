package client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/internal/logging/central"
	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	errPingFormat          = "failed to ping mongodb: %w"
	errProbeFormat         = "failed to probe mongodb: %w"
	errCreateClientFormat  = "failed to create mongodb client: %w"
	errInvalidConfigFormat = "invalid mongodb config: %w"
)

// attemptConnection establishes a driver client with specified TLS and verifies ping and probe.
func attemptConnection(
	ctx context.Context,
	cfg *config.Config,
	enableTLS bool,
	opts ...Option,
) (*Client, error) {
	optionsContainer := NewOptions(opts...)
	central.DefaultWorker().Logger().InfoContext(ctx, "resolving mongodb connection options",
		"host", cfg.Host,
		"port", cfg.Port,
		"auth_source", cfg.AuthSource,
		"tls", enableTLS,
		"uuid_representation", cfg.UUIDRepresentation,
	)
	clientOptions := BuildClientOptionsWithTLS(cfg, enableTLS, optionsContainer.DriverOptions...)

	rawClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf(errCreateClientFormat, err)
	}

	if err := verifyPingAndProbe(ctx, cfg, rawClient); err != nil {
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

	return nil, handleConnectionError(ctx, err)
}

// Connect loads MongoDB configuration from environment variables
// (OHM9996_MONGODB_*), establishes a connection, verifies ping connectivity
// and random collection/doc probe, and returns a managed Client.
func Connect(ctx context.Context, opts ...Option) (*Client, error) {
	fn := callstack.DecorateContextFunc("Connect", func(ctx context.Context) (*Client, error) {
		cfg, err := config.LoadConfig()
		if err != nil {
			return nil, err
		}
		return ConnectWithConfig(ctx, cfg, opts...)
	})
	return fn(ctx)
}

func validateAndLogConfig(ctx context.Context, cfg *config.Config) error {
	if cfg == nil {
		return config.ErrNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf(errInvalidConfigFormat, err)
	}
	central.DefaultWorker().Logger().InfoContext(ctx, "initializing mongodb client",
		"host", cfg.Host,
		"port", cfg.Port,
		"protocol", cfg.Protocol,
		"database", cfg.Database,
	)
	return nil
}

// ConnectWithConfig connects to MongoDB using the provided Config and Options.
// It verifies mandatory ping and probe connectivity using the provided context.
func ConnectWithConfig(
	ctx context.Context,
	cfg *config.Config,
	opts ...Option,
) (*Client, error) {
	fn := callstack.DecorateContextFunc("ConnectWithConfig", func(ctx context.Context) (*Client, error) {
		if err := validateAndLogConfig(ctx, cfg); err != nil {
			return nil, err
		}
		client, err := initAndPingClient(ctx, cfg, opts...)
		if err != nil {
			return nil, err
		}
		central.DefaultWorker().Logger().InfoContext(ctx, "mongodb client initialization complete",
			"database", resolveDatabaseName(cfg),
		)
		return client, nil
	})
	return fn(ctx)
}
