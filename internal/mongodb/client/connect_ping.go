package client

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

var (
	pingClient  = defaultPingClient
	probeClient = defaultProbeClient
)

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, nil)
}

// resolveDatabaseName extracts the target database name from config.
func resolveDatabaseName(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	return cfg.Database
}

func verifyPing(ctx context.Context, rawClient *mongo.Client) error {
	if err := pingClient(ctx, rawClient); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(fmt.Errorf(errPingFormat, err), disconnectErr)
		}
		return fmt.Errorf(errPingFormat, err)
	}
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, "mongodb ping succeeded")
	return nil
}

func verifyProbe(ctx context.Context, rawClient *mongo.Client, dbName string) error {
	if err := probeClient(ctx, rawClient, dbName); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(fmt.Errorf(errProbeFormat, err), disconnectErr)
		}
		return fmt.Errorf(errProbeFormat, err)
	}
	return nil
}

// verifyPingAndProbe verifies mandatory ping and random collection/document probe,
// disconnecting the client if either check fails.
func verifyPingAndProbe(
	ctx context.Context,
	cfg *config.Config,
	rawClient *mongo.Client,
) error {
	if err := verifyPing(ctx, rawClient); err != nil {
		return err
	}
	probeCtx := withProbeTimeout(ctx, resolveProbeTimeout(cfg))
	return verifyProbe(probeCtx, rawClient, resolveDatabaseName(cfg))
}
