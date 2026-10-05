package client

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

var (
	pingClient          = defaultPingClient
	pingFirestoreClient = defaultPingFirestoreClient
	probeClient         = defaultProbeClient
)

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, nil)
}

func defaultPingFirestoreClient(ctx context.Context, rawClient *mongo.Client, dbName string) error {
	db := rawClient.Database(resolveDBName(dbName))
	res := db.Collection("ping").FindOne(ctx, bson.D{})
	if err := res.Err(); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	return nil
}

// pingFirestore performs a dummy read on the target database/collection to verify connectivity
// when connecting to Google Cloud Firestore with MongoDB compatibility.
func pingFirestore(ctx context.Context, rawClient *mongo.Client, dbName string) error {
	return pingFirestoreClient(ctx, rawClient, dbName)
}

// resolveDatabaseName extracts the target database name from config.
func resolveDatabaseName(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	return cfg.Database
}

func dispatchPing(ctx context.Context, cfg *config.Config, rawClient *mongo.Client) error {
	if cfg.IsFirestore() {
		return pingFirestore(ctx, rawClient, resolveDatabaseName(cfg))
	}
	return pingClient(ctx, rawClient)
}

func handlePingDisconnect(ctx context.Context, rawClient *mongo.Client, err error) error {
	if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
		return errors.Join(fmt.Errorf(errPingFormat, err), disconnectErr)
	}
	return fmt.Errorf(errPingFormat, err)
}

func verifyPing(ctx context.Context, cfg *config.Config, rawClient *mongo.Client) error {
	if err := dispatchPing(ctx, cfg, rawClient); err != nil {
		return handlePingDisconnect(ctx, rawClient, err)
	}
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, "mongodb ping succeeded")
	return nil
}

func verifyProbe(ctx context.Context, rawClient *mongo.Client, dbName string, isFirestore ...bool) error {
	managedClient := NewClient(rawClient, dbName, isFirestore...)
	if err := probeClient(ctx, managedClient, dbName); err != nil {
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
	if err := verifyPing(ctx, cfg, rawClient); err != nil {
		return err
	}
	probeCtx := withProbeTimeout(ctx, resolveProbeTimeout(cfg))
	return verifyProbe(probeCtx, rawClient, resolveDatabaseName(cfg), cfg.IsFirestore())
}
