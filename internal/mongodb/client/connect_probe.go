package client

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
)

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

func resolveProbeDuration(d, defaultDur time.Duration) time.Duration {
	if d <= 0 {
		return defaultDur
	}
	return d
}

func probeListCollectionsWithoutMaxTime(
	ctx context.Context,
	db *mongo.Database,
	timeout time.Duration,
) ([]string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, resolveProbeDuration(timeout, defaultProbeTimeout))
	defer cancel()

	colls, err := db.ListCollectionNames(timeoutCtx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections without maxTimeMS: %w", err)
	}
	return colls, nil
}

func probeListCollectionsWithMaxTime(
	ctx context.Context,
	db *mongo.Database,
	timeout time.Duration,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, resolveProbeDuration(timeout, defaultProbeTimeout))
	defer cancel()
	cmd := bson.D{{Key: "listCollections", Value: 1}}
	if err := db.RunCommand(timeoutCtx, cmd).Err(); err != nil {
		return fmt.Errorf("failed to list collections with maxTimeMS: %w", err)
	}
	return nil
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
	timeout := probeTimeoutFromContext(ctx)
	db := rawClient.Database(resolveDBName(dbName))
	colls, err := probeListCollectionsWithoutMaxTime(ctx, db, timeout)
	if err != nil {
		return err
	}
	if err := probeListCollectionsWithMaxTime(ctx, db, timeout); err != nil {
		return err
	}
	return probeDocument(ctx, db, chooseTargetColl(colls))
}
