package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

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
