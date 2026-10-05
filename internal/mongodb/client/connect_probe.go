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
	return "probe"
}

func resolveProbeDuration(d, defaultDur time.Duration) time.Duration {
	if d <= 0 {
		return defaultDur
	}
	return d
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

func probeCollections(
	ctx context.Context,
	c *Client,
	db *mongo.Database,
	timeout time.Duration,
) ([]string, error) {
	if c.IsFirestore() {
		return probeFirestoreCollections(ctx, db)
	}
	return probeStandardCollections(ctx, db, timeout)
}

func resolveProbeDatabase(c *Client, dbName string) (*mongo.Database, error) {
	if c == nil {
		return nil, ErrNilClient
	}
	db := c.Database(resolveDBName(dbName))
	if db == nil {
		return nil, ErrNilClient
	}
	return db, nil
}

func defaultProbeClient(ctx context.Context, c *Client, dbName string) error {
	db, err := resolveProbeDatabase(c, dbName)
	if err != nil {
		return err
	}
	colls, err := probeCollections(ctx, c, db, probeTimeoutFromContext(ctx))
	if err != nil {
		return err
	}
	return probeDocument(ctx, db, chooseTargetColl(colls))
}
