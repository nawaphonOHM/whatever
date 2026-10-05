package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
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
	pingClient          = defaultPingClient
	pingFirestoreClient = defaultPingFirestoreClient
	probeClient         = defaultProbeClient
)

func defaultPingClient(ctx context.Context, rawClient *mongo.Client) error {
	return rawClient.Ping(ctx, readpref.Primary())
}

func defaultPingFirestoreClient(ctx context.Context, rawClient *mongo.Client, dbName string) error {
	db := rawClient.Database(resolveDBName(dbName))
	res := db.Collection("ping").FindOne(ctx, bson.D{})
	if err := res.Err(); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	return nil
}

func pingFirestore(ctx context.Context, rawClient *mongo.Client, dbName string) error {
	return pingFirestoreClient(ctx, rawClient, dbName)
}

// resolveDatabaseName extracts the target database name from options.
func resolveDatabaseName(o *Options) string {
	if o == nil {
		return ""
	}
	return o.Database
}

func dispatchPing(ctx context.Context, o *Options, rawClient *mongo.Client) error {
	if o != nil && o.IsFirestore() {
		return pingFirestore(ctx, rawClient, resolveDatabaseName(o))
	}
	return pingClient(ctx, rawClient)
}

func verifyClientPing(ctx context.Context, o *Options, rawClient *mongo.Client) error {
	if err := dispatchPing(ctx, o, rawClient); err != nil {
		if disconnectErr := rawClient.Disconnect(ctx); disconnectErr != nil {
			return errors.Join(fmt.Errorf(errPingFormat, err), disconnectErr)
		}
		return fmt.Errorf(errPingFormat, err)
	}
	return nil
}

func verifyClientProbe(ctx context.Context, rawClient *mongo.Client, database string, isFirestore ...bool) error {
	managedClient := client.NewClient(rawClient, database, isFirestore...)
	if err := probeClient(ctx, managedClient, database); err != nil {
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
	if err := verifyClientPing(ctx, o, rawClient); err != nil {
		return err
	}
	probeCtx := withProbeTimeout(ctx, resolveProbeTimeout(o))
	return verifyClientProbe(probeCtx, rawClient, resolveDatabaseName(o), o.IsFirestore())
}
