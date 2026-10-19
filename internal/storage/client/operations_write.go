// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func (c *Client) buildPutObjectInput(input contracts.PutObjectInput, body io.Reader) *s3.PutObjectInput {
	params := &s3.PutObjectInput{
		Bucket: aws.String(c.resolveBucket(input.Bucket)),
		Key:    aws.String(input.Key),
		Body:   body,
	}
	if input.ContentType != "" {
		params.ContentType = aws.String(input.ContentType)
	}
	if len(input.Metadata) > 0 {
		params.Metadata = input.Metadata
	}
	return params
}

// putSingle sends one PutObject request. The body must be nil or seekable so the SDK can sign the
// payload and determine Content-Length on both HTTP and HTTPS endpoints.
func (c *Client) putSingle(
	ctx context.Context,
	input contracts.PutObjectInput,
	body io.Reader,
) (*contracts.PutObjectOutput, error) {
	out, err := c.rawClient.PutObject(ctx, c.buildPutObjectInput(input, body))
	if err != nil {
		return nil, err
	}
	return &contracts.PutObjectOutput{
		ETag:      aws.ToString(out.ETag),
		VersionID: aws.ToString(out.VersionId),
	}, nil
}

func isSmallSeekable(body io.Reader) bool {
	return isSeekable(body) && !isLargeSeekable(body)
}

func isDirectBody(body io.Reader) bool {
	return body == nil || isSmallSeekable(body)
}

func (c *Client) putBody(ctx context.Context, input contracts.PutObjectInput) (*contracts.PutObjectOutput, error) {
	if isOversizedSeekable(input.Body) {
		return nil, errTooManyParts
	}
	if isDirectBody(input.Body) {
		return c.putSingle(ctx, input, input.Body)
	}
	return c.putStream(ctx, input)
}

// PutObject uploads an object to the storage bucket from any io.Reader.
//
// Seekable bodies (bytes.Reader, strings.Reader, *os.File, ...) are sent in a single signed request.
// Other readers (pipes, network streams, bytes.Buffer, unknown length) are buffered in bounded memory:
// bodies up to one part are uploaded in a single request and larger ones use a multipart upload that is
// aborted on failure (the maximum size is 10000 parts of 8 MiB). RequestTimeout bounds the whole upload,
// so configure it generously for large objects or use pre-signed URLs for direct transfers.
func (c *Client) PutObject(ctx context.Context, input contracts.PutObjectInput) (*contracts.PutObjectOutput, error) {
	if c.isNil() {
		return nil, contracts.ErrNilClient
	}
	reqCtx, cancel := c.requestContext(ctx)
	defer cancel()
	out, err := c.putBody(reqCtx, input)
	if err != nil {
		return nil, mapError(err)
	}
	return out, nil
}

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
