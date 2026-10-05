package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func probeListCollectionsWithoutMaxTime(
	ctx context.Context,
	db *mongo.Database,
	timeout time.Duration,
) ([]string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, resolveProbeDuration(timeout, DefaultSocketTimeout))
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
	timeoutCtx, cancel := context.WithTimeout(ctx, resolveProbeDuration(timeout, DefaultSocketTimeout))
	defer cancel()
	cmd := bson.D{{Key: "listCollections", Value: 1}}
	if err := db.RunCommand(timeoutCtx, cmd).Err(); err != nil {
		return fmt.Errorf("failed to list collections with maxTimeMS: %w", err)
	}
	return nil
}

func probeStandardCollections(
	ctx context.Context,
	db *mongo.Database,
	timeout time.Duration,
) ([]string, error) {
	colls, err := probeListCollectionsWithoutMaxTime(ctx, db, timeout)
	if err != nil {
		return nil, err
	}
	if err := probeListCollectionsWithMaxTime(ctx, db, timeout); err != nil {
		return nil, err
	}
	return colls, nil
}
