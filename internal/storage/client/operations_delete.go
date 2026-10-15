// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// DeleteObject removes an object by key from the storage bucket.
func (c *Client) DeleteObject(
	ctx context.Context,
	input contracts.DeleteObjectInput,
) (*contracts.DeleteObjectOutput, error) {
	if c.isNil() {
		return nil, contracts.ErrNilClient
	}
	reqCtx, cancel := c.requestContext(ctx)
	defer cancel()
	bucket := c.resolveBucket(input.Bucket)
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(input.Key),
	}
	out, err := c.rawClient.DeleteObject(reqCtx, params)
	if err != nil {
		return nil, mapError(err)
	}
	return &contracts.DeleteObjectOutput{
		VersionID: aws.ToString(out.VersionId),
	}, nil
}
