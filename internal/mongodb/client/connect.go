package client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	_ "github.com/nawaphonOHM/whatever/v3/internal/logging/core" // side-effect: eagerly initialize structured logging
	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

const (
	errPingFormat          = "failed to ping mongodb: %w"
	errProbeFormat         = "failed to probe mongodb: %w"
	errCreateClientFormat  = "failed to create mongodb client: %w"
	errInvalidConfigFormat = "invalid mongodb config: %w"
)

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
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
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
		// False positive: borrowing the central worker singleton for logging is leak-free and safe;
		// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
		central.DefaultWorker().Logger().InfoContext(ctx, "mongodb client initialization complete",
			"database", resolveDatabaseName(cfg),
		)
		return client, nil
	})
	return fn(ctx)
}

// SetMockPing overrides pingClient for testing and returns a restore function.
func SetMockPing(fn func(context.Context, *mongo.Client) error) func() {
	prev := pingClient
	if fn != nil {
		pingClient = fn
	} else {
		pingClient = defaultPingClient
	}
	return func() {
		pingClient = prev
	}
}

// SetMockPingFirestore overrides pingFirestoreClient for testing and returns a restore function.
func SetMockPingFirestore(fn func(context.Context, *mongo.Client, string) error) func() {
	prev := pingFirestoreClient
	if fn != nil {
		pingFirestoreClient = fn
	} else {
		pingFirestoreClient = defaultPingFirestoreClient
	}
	return func() {
		pingFirestoreClient = prev
	}
}

// SetMockProbe overrides probeClient for testing and returns a restore function.
func SetMockProbe(fn func(context.Context, *Client, string) error) func() {
	prev := probeClient
	if fn != nil {
		probeClient = fn
	} else {
		probeClient = defaultProbeClient
	}
	return func() {
		probeClient = prev
	}
}
