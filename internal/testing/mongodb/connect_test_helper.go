package mongodb

import (
	"context"
	"testing"
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
