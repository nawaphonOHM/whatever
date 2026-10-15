// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// Ping probes connectivity and accessibility of the default bucket. The probe is bounded by the
// configured PingTimeout. It issues HeadBucket and, when that is denied (some providers require extra
// bucket-level permission), falls back to a ListObjectsV2 request limited to a single key.
func (c *Client) Ping(ctx context.Context) error {
	if c.isNil() {
		return contracts.ErrNilClient
	}
	return c.probeBucket(ctx)
}

func (c *Client) headBucket(ctx context.Context, bucket string) error {
	_, err := c.rawClient.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	return mapBucketError(err)
}

func (c *Client) listProbe(ctx context.Context, bucket string) error {
	_, err := c.rawClient.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int32(1),
	})
	return mapBucketError(err)
}

func (c *Client) reachBucket(ctx context.Context, bucket string) error {
	err := c.headBucket(ctx, bucket)
	if errors.Is(err, contracts.ErrAccessDenied) {
		return c.listProbe(ctx, bucket)
	}
	return err
}

func (c *Client) probeBucket(ctx context.Context) error {
	if c.defaultBucket == "" {
		return contracts.ErrBucketNotFound
	}
	pingCtx, cancel := withTimeout(ctx, c.pingTimeout)
	defer cancel()
	return c.reachBucket(pingCtx, c.defaultBucket)
}
