package mongodb

import (
	"context"
	"testing"

	intclient "github.com/nawaphonOHM/whatever/v3/internal/mongodb/client"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// registerCleanup attaches client.Close to tb.Cleanup.
func registerCleanup(tb testing.TB, client *TestClient) {
	tb.Cleanup(func() {
		if err := client.Close(); err != nil {
			tb.Logf("failed to close mongodb test client: %v", err)
		}
	})
}

// NewTestClient creates a TestClient using functional options and registers automatic cleanup on tb.
// If connection fails, tb.Fatalf is called to fail the test immediately.
func NewTestClient(tb testing.TB, opts ...Option) *TestClient {
	if tb == nil {
		panic(ErrNilTestingTB)
	}
	tb.Helper()
	client, err := Connect(context.Background(), opts...)
	if err != nil {
		tb.Fatalf("failed to connect to mongodb: %v", err)
		return nil
	}
	registerCleanup(tb, client)
	return client
}

// NewTestClientURI creates a TestClient using a URI string and registers automatic cleanup on tb.
// If connection fails, tb.Fatalf is called to fail the test immediately.
func NewTestClientURI(tb testing.TB, uri string, opts ...Option) *TestClient {
	if tb == nil {
		panic(ErrNilTestingTB)
	}
	tb.Helper()
	client, err := ConnectURI(context.Background(), uri, opts...)
	if err != nil {
		tb.Fatalf("failed to connect to mongodb uri %q: %v", uri, err)
		return nil
	}
	registerCleanup(tb, client)
	return client
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

// SetMockPingFirestore overrides pingFirestoreClient for testing and returns a restore function.
func SetMockPingFirestore(fn func(context.Context, *mongo.Client, string) error) func() {
	prev := pingFirestoreClient
	if fn != nil {
		pingFirestoreClient = fn
	} else {
		pingFirestoreClient = defaultPingFirestoreClient
	}
	return func() {
		pingFirestoreClient = prev
	}
}

// SetMockProbe overrides probeClient for testing and returns a restore function.
func SetMockProbe(fn func(context.Context, *intclient.Client, string) error) func() {
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
