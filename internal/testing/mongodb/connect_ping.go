package mongodb

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	errPingFormat         = "failed to ping mongodb: %w"
	errProbeFormat        = "failed to probe mongodb: %w"
	errCreateClientFormat = "failed to create mongodb client: %w"
)

var (
	pingClient  = defaultPingClient
	probeClient = defaultProbeClient
)

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, readpref.Primary())
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

// SetMockPing overrides pingClient for testing and returns a restore function.
func SetMockPing(fn func(context.Context, *mongo.Client) error) func() {
	prev := pingClient
	if fn != nil {
		pingClient = fn
	} else {
		pingClient = defaultPingClient
	}
	return func() {
		pingClient = prev
	}
}

// SetMockProbe overrides probeClient for testing and returns a restore function.
func SetMockProbe(fn func(context.Context, *mongo.Client, string) error) func() {
	prev := probeClient
	if fn != nil {
		probeClient = fn
	} else {
		probeClient = defaultProbeClient
	}
	return func() {
		probeClient = prev
	}
}

func verifyClientPing(ctx context.Context, rawClient *mongo.Client) error {
	if err := pingClient(ctx, rawClient); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(fmt.Errorf(errPingFormat, err), disconnectErr)
		}
		return fmt.Errorf(errPingFormat, err)
	}
	return nil
}

func verifyClientProbe(ctx context.Context, rawClient *mongo.Client, database string) error {
	if err := probeClient(ctx, rawClient, database); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(fmt.Errorf(errProbeFormat, err), disconnectErr)
		}
		return fmt.Errorf(errProbeFormat, err)
	}
	return nil
}

// verifyClientPingAndProbe verifies mandatory ping and collection/document probe,
// disconnecting the client if either fails.
func verifyClientPingAndProbe(ctx context.Context, o *Options, rawClient *mongo.Client) error {
	if err := verifyClientPing(ctx, rawClient); err != nil {
		return err
	}
	return verifyClientProbe(ctx, rawClient, o.Database)
}
