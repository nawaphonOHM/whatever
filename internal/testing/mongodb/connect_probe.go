package mongodb

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func probeListCollectionsWithoutMaxTime(
	ctx context.Context,
	db *mongo.Database,
	timeout time.Duration,
) ([]string, error) {
	effectiveTimeout := timeout
	if effectiveTimeout <= 0 {
		effectiveTimeout = DefaultSocketTimeout
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, effectiveTimeout)
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
	effectiveTimeout := timeout
	if effectiveTimeout <= 0 {
		effectiveTimeout = DefaultSocketTimeout
	}
	cmd := bson.D{
		{Key: "listCollections", Value: 1},
		{Key: "maxTimeMS", Value: effectiveTimeout.Milliseconds()},
	}
	if err := db.RunCommand(ctx, cmd).Err(); err != nil {
		return fmt.Errorf("failed to list collections with maxTimeMS: %w", err)
	}
	return nil
}

func probeDocument(ctx context.Context, db *mongo.Database, coll string) error {
	res := db.Collection(coll).FindOne(ctx, bson.D{})
	if err := res.Err(); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to probe document in collection %q: %w", coll, err)
	}
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
