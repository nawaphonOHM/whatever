package client

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/internal/logging/central"
	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

var (
	pingClient  = defaultPingClient
	probeClient = defaultProbeClient
)

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, nil)
}

func resolveDBName(dbName string) string {
	if dbName == "" {
		return "admin"
	}
	return dbName
}

func chooseTargetColl(colls []string) string {
	if len(colls) > 0 {
		return colls[rand.IntN(len(colls))]
	}
	return "__probe__"
}

func probeDocument(ctx context.Context, db *mongo.Database, coll string) error {
	res := db.Collection(coll).FindOne(ctx, bson.D{})
	if err := res.Err(); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to probe document in collection %q: %w", coll, err)
	}
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, "mongodb probe succeeded", "target_collection", coll)
	return nil
}

func defaultProbeClient(ctx context.Context, rawClient *mongo.Client, dbName string) error {
	db := rawClient.Database(resolveDBName(dbName))
	colls, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to list collections for probe: %w", err)
	}
	return probeDocument(ctx, db, chooseTargetColl(colls))
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
	return verifyProbe(ctx, rawClient, resolveDatabaseName(cfg))
}
