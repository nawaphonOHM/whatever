package client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

func logConnectionOptions(ctx context.Context, cfg *config.Config, enableTLS bool) {
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, "resolving mongodb connection options",
		"host", cfg.Host,
		"port", cfg.Port,
		"auth_source", cfg.AuthSource,
		"tls", enableTLS,
		"uuid_representation", cfg.UUIDRepresentation,
	)
}

// attemptConnection establishes a driver client with specified TLS and verifies ping and probe.
func attemptConnection(
	ctx context.Context,
	cfg *config.Config,
	enableTLS bool,
	opts ...Option,
) (*Client, error) {
	optionsContainer := NewOptions(opts...)
	logConnectionOptions(ctx, cfg, enableTLS)
	clientOptions := BuildClientOptionsWithTLS(cfg, enableTLS, optionsContainer.DriverOptions...)

	rawClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf(errCreateClientFormat, err)
	}

	if err := verifyPingAndProbe(ctx, cfg, rawClient); err != nil {
		return nil, err
	}
	return NewClient(rawClient, resolveDatabaseName(cfg), cfg.IsFirestore()), nil
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
