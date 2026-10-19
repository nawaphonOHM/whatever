// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func mapGetObjectOutput(out *s3.GetObjectOutput) *contracts.GetObjectOutput {
	return &contracts.GetObjectOutput{
		Body:          out.Body,
		ContentLength: aws.ToInt64(out.ContentLength),
		ContentType:   aws.ToString(out.ContentType),
		ETag:          aws.ToString(out.ETag),
		Metadata:      out.Metadata,
		LastModified:  aws.ToTime(out.LastModified),
	}
}

// GetObject retrieves an object stream and metadata from the storage bucket.
//
// Ownership: the caller owns GetObjectOutput.Body and must Close it exactly when done; Close releases
// the connection and the request context. RequestTimeout only bounds the wait for the response headers,
// so large bodies may stream for as long as the caller's ctx allows (cancel ctx to abort reading).
func (c *Client) GetObject(ctx context.Context, input contracts.GetObjectInput) (*contracts.GetObjectOutput, error) {
	if c.isNil() {
		return nil, contracts.ErrNilClient
	}
	reqCtx, cancel := context.WithCancelCause(ctx)
	guard := startHeaderTimer(c.requestTimeout, cancel)
	out, err := c.rawClient.GetObject(reqCtx, &s3.GetObjectInput{
		Bucket: aws.String(c.resolveBucket(input.Bucket)),
		Key:    aws.String(input.Key),
	})
	return finishGet(getResult{out: out, err: err, cancel: cancel, guard: guard})
}

func mapHeadObjectOutput(out *s3.HeadObjectOutput) *contracts.HeadObjectOutput {
	return &contracts.HeadObjectOutput{
		ContentLength: aws.ToInt64(out.ContentLength),
		ContentType:   aws.ToString(out.ContentType),
		ETag:          aws.ToString(out.ETag),
		Metadata:      out.Metadata,
		LastModified:  aws.ToTime(out.LastModified),
	}
}

// mapHeadError maps a HEAD failure. A body-less 404 cannot tell a missing object from a missing
// bucket, so on object-not-found the bucket is probed to return ErrBucketNotFound when appropriate.
func (c *Client) mapHeadError(ctx context.Context, bucket string, err error) error {
	mapped := mapError(err)
	if !errors.Is(mapped, contracts.ErrObjectNotFound) {
		return mapped
	}
	if errors.Is(c.headBucket(ctx, bucket), contracts.ErrBucketNotFound) {
		return fmt.Errorf(errWrapFmt, contracts.ErrBucketNotFound, err)
	}
	return mapped
}

// HeadObject retrieves object metadata without downloading the object payload.
func (c *Client) HeadObject(ctx context.Context, input contracts.HeadObjectInput) (*contracts.HeadObjectOutput, error) {
	if c.isNil() {
		return nil, contracts.ErrNilClient
	}
	reqCtx, cancel := c.requestContext(ctx)
	defer cancel()
	bucket := c.resolveBucket(input.Bucket)
	out, err := c.rawClient.HeadObject(reqCtx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(input.Key),
	})
	if err != nil {
		return nil, c.mapHeadError(reqCtx, bucket, err)
	}
	return mapHeadObjectOutput(out), nil
}

func mapExistsResult(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if errors.Is(err, contracts.ErrObjectNotFound) {
		return false, nil
	}
	return false, err
}

// Exists checks whether a specific object key exists in the bucket. A missing object yields
// (false, nil); a missing bucket, access denial (403) and transport failures are returned as errors.
func (c *Client) Exists(ctx context.Context, bucket, key string) (bool, error) {
	if c.isNil() {
		return false, contracts.ErrNilClient
	}
	b := c.resolveBucket(bucket)
	_, err := c.HeadObject(ctx, contracts.HeadObjectInput{Bucket: b, Key: key})
	return mapExistsResult(err)
}
