package client

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

const defaultProbeTimeout = 10 * time.Second

type probeTimeoutKey struct{}

func withProbeTimeout(ctx context.Context, timeout time.Duration) context.Context {
	if timeout <= 0 {
		return ctx
	}
	return context.WithValue(ctx, probeTimeoutKey{}, timeout)
}

func probeTimeoutFromContext(ctx context.Context) time.Duration {
	if v, ok := ctx.Value(probeTimeoutKey{}).(time.Duration); ok && v > 0 {
		return v
	}
	return defaultProbeTimeout
}

func resolveProbeTimeout(cfg *config.Config) time.Duration {
	if cfg == nil || cfg.SocketTimeout <= 0 {
		return defaultProbeTimeout
	}
	return cfg.SocketTimeout
}

func probeFirestoreCollections(
	ctx context.Context,
	db *mongo.Database,
) ([]string, error) {
	colls, err := db.ListCollectionNames(context.WithoutCancel(ctx), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections without maxTimeMS: %w", err)
	}
	return colls, nil
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
