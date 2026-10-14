// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func mapPresignedURL(req *signer.PresignedHTTPRequest, expiry time.Duration) *contracts.PresignedURL {
	headers := make(map[string][]string, len(req.SignedHeader))
	for k, v := range req.SignedHeader {
		headers[k] = append([]string(nil), v...)
	}
	return &contracts.PresignedURL{
		Headers:   headers,
		URL:       req.URL,
		Method:    req.Method,
		ExpiresIn: expiry,
	}
}

func buildPresignPutInput(bucket string, input contracts.PresignPutObjectInput) *s3.PutObjectInput {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(input.Key),
	}
	if input.ContentType != "" {
		params.ContentType = aws.String(input.ContentType)
	}
	return params
}

// PresignGetObject generates a pre-signed URL for downloading an object.
// ExpiresIn of 0 uses the configured default; negative values or values above 7 days
// are rejected with ErrInvalidExpiry.
func (c *Client) PresignGetObject(
	ctx context.Context,
	input contracts.PresignGetObjectInput,
) (*contracts.PresignedURL, error) {
	if c.isPresignNil() {
		return nil, contracts.ErrNilClient
	}
	expiry, err := c.resolvePresignExpiry(input.ExpiresIn)
	if err != nil {
		return nil, err
	}
	params := &s3.GetObjectInput{
		Bucket: aws.String(c.resolveBucket(input.Bucket)),
		Key:    aws.String(input.Key),
	}
	return c.presign(ctx, expiry, func(reqCtx context.Context) (*signer.PresignedHTTPRequest, error) {
		return c.presignClient.PresignGetObject(reqCtx, params, s3.WithPresignExpires(expiry))
	})
}

// PresignPutObject generates a pre-signed URL for uploading an object. When ContentType is set it is
// signed: the uploader must send the returned Headers (notably Content-Type) unchanged.
func (c *Client) PresignPutObject(
	ctx context.Context,
	input contracts.PresignPutObjectInput,
) (*contracts.PresignedURL, error) {
	if c.isPresignNil() {
		return nil, contracts.ErrNilClient
	}
	expiry, err := c.resolvePresignExpiry(input.ExpiresIn)
	if err != nil {
		return nil, err
	}
	params := buildPresignPutInput(c.resolveBucket(input.Bucket), input)
	return c.presign(ctx, expiry, func(reqCtx context.Context) (*signer.PresignedHTTPRequest, error) {
		return c.presignClient.PresignPutObject(reqCtx, params, s3.WithPresignExpires(expiry))
	})
}

func (c *Client) presign(
	ctx context.Context,
	expiry time.Duration,
	sign func(context.Context) (*signer.PresignedHTTPRequest, error),
) (*contracts.PresignedURL, error) {
	reqCtx, cancel := c.requestContext(ctx)
	defer cancel()
	req, err := sign(reqCtx)
	if err != nil {
		return nil, mapError(err)
	}
	return mapPresignedURL(req, expiry), nil
}
