package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// attemptConnection establishes a driver client and optionally verifies ping.
func attemptConnection(
	ctx context.Context,
	o *Options,
) (*TestClient, error) {
	clientOptions := BuildClientOptions(o)
	rawClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf(errCreateClientFormat, err)
	}

	if err := verifyClientPing(ctx, o, rawClient); err != nil {
		return nil, err
	}
	return NewClient(rawClient, o.Database), nil
}

// fallbackTLSAttempt retries connection with TLS enabled if error indicates TLS requirement.
func fallbackTLSAttempt(
	ctx context.Context,
	o *Options,
	err error,
) (*TestClient, error) {
	if !o.EnableTLS && isTLSError(err) {
		tlsOptions := *o
		tlsOptions.EnableTLS = true
		return attemptConnection(ctx, &tlsOptions)
	}
	return nil, err
}

// connectWithOptions handles connection with two-phase TLS fallback.
func connectWithOptions(ctx context.Context, o *Options) (*TestClient, error) {
	if err := validateOptions(o); err != nil {
		return nil, err
	}

	client, err := attemptConnection(ctx, o)
	if err != nil {
		return fallbackTLSAttempt(ctx, o, err)
	}
	return client, nil
}

// Connect creates a TestClient using functional options with two-phase TLS fallback.
func Connect(ctx context.Context, opts ...Option) (*TestClient, error) {
	o := NewOptions(opts...)
	return connectWithOptions(ctx, o)
}

// ConnectURI connects to MongoDB using a connection URI string and functional options.
func ConnectURI(ctx context.Context, uri string, opts ...Option) (*TestClient, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, ErrEmptyURI
	}
	allOpts := append([]Option{WithURI(uri)}, opts...)
	o := NewOptions(allOpts...)
	if o.Database == "" {
		o.Database = extractDatabaseFromURI(uri)
	}
	return connectWithOptions(ctx, o)
}
