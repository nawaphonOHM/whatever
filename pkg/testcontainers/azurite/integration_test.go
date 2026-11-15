//go:build testcontainers

package azurite_test

import (
	"context"
	"testing"
)

func TestS3BucketAndObjectLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), containerTimeout)
	defer cancel()

	endpoint, _, _ := startS3GatewayFixture(t, ctx)

	names := newS3ResourceNames()
	client, markBucketDeleted := newAzuriteStorageClient(t, ctx, endpoint, names)
	(s3LifecycleFixture{
		ctx:    ctx,
		client: client,
		names:  names,
	}).verify(t, markBucketDeleted)
}
