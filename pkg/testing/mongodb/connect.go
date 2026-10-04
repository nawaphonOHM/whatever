package mongodb

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/testing/mongodb"
)

// SetMockPing overrides the mock ping function for testing and returns a restore function.
func SetMockPing(fn func(context.Context, *mongo.Client) error) func() {
	return mongodb.SetMockPing(fn)
}

// SetMockProbe overrides the mock probe function for testing and returns a restore function.
func SetMockProbe(fn func(context.Context, *mongo.Client, string) error) func() {
	return mongodb.SetMockProbe(fn)
}

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
