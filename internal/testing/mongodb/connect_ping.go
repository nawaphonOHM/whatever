package mongodb

import (
	"context"
	"errors"
	"fmt"

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

// resolveDatabaseName extracts the target database name from options.
func resolveDatabaseName(o *Options) string {
	if o == nil {
		return ""
	}
	return o.Database
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
	probeCtx := withProbeTimeout(ctx, resolveProbeTimeout(o))
	return verifyClientProbe(probeCtx, rawClient, resolveDatabaseName(o))
}
