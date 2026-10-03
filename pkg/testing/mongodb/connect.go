package mongodb

import (
	"context"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/testing/mongodb"
)

// Connect creates a TestClient using functional options without process exit behavior.
func Connect(ctx context.Context, opts ...Option) (*TestClient, error) {
	return mongodb.Connect(ctx, opts...)
}

// ConnectURI connects to MongoDB using a connection URI string and functional options.
func ConnectURI(ctx context.Context, uri string, opts ...Option) (*TestClient, error) {
	return mongodb.ConnectURI(ctx, uri, opts...)
}

// NewTestClient creates a TestClient using functional options, asserting connectivity and registering tb.Cleanup.
func NewTestClient(tb testing.TB, opts ...Option) *TestClient {
	return mongodb.NewTestClient(tb, opts...)
}

// NewTestClientURI creates a TestClient using a URI string, asserting connectivity and registering tb.Cleanup.
func NewTestClientURI(tb testing.TB, uri string, opts ...Option) *TestClient {
	return mongodb.NewTestClientURI(tb, uri, opts...)
}
