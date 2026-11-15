//go:build testcontainers

package azurite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func TestS3GatewayCreatesBucketInAzurite(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), containerTimeout)
	defer cancel()
	endpoint, client, _ := startS3GatewayFixture(t, ctx)
	require.NotEmpty(t, endpoint)

	bucket := fmt.Sprintf("azurite-s3-gateway-%d", time.Now().UnixNano())
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
	require.NoError(t, err, "authenticated S3 bucket creation must reach Azurite through S3Proxy")
	registerS3BucketCleanup(t, client, bucket)
}

func TestS3GatewayRequestsFailWhenAzuriteIsUnavailable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), containerTimeout)
	defer cancel()
	_, client, stopAzurite := startS3GatewayFixture(t, ctx)
	require.NoError(t, stopAzurite(ctx))

	bucket := fmt.Sprintf("azurite-s3-outage-%d", time.Now().UnixNano())
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
	if err == nil {
		registerS3BucketCleanup(t, client, bucket)
	}
	require.Error(t, err, "S3Proxy must fail bucket creation when its Azurite backend is unavailable")
}

func registerS3BucketCleanup(t *testing.T, client *s3.Client, bucket string) func() {
	t.Helper()
	bucketDeleted := false
	t.Cleanup(func() {
		if bucketDeleted {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), s3GatewayCleanupLimit)
		defer cancel()
		if _, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)}); err != nil {
			t.Errorf("delete S3 gateway test bucket: %v", err)
		}
	})
	return func() {
		bucketDeleted = true
	}
}
