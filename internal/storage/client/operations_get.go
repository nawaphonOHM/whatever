// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

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
