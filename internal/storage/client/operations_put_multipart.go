// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const maxMultipartParts = 10000

var errTooManyParts = errors.New("storage upload exceeds the maximum number of multipart parts")

type multipartUpload struct {
	client   *Client
	bucket   *string
	key      *string
	uploadID *string
	parts    []types.CompletedPart
}

func (m *multipartUpload) checksumAlgorithm() types.ChecksumAlgorithm {
	if m.client.provider == contracts.ProviderAWS {
		return types.ChecksumAlgorithmCrc32
	}
	return ""
}

func (m *multipartUpload) create(ctx context.Context, input contracts.PutObjectInput) error {
	params := &s3.CreateMultipartUploadInput{
		Bucket: m.bucket, Key: m.key, Metadata: input.Metadata, ChecksumAlgorithm: m.checksumAlgorithm(),
	}
	if input.ContentType != "" {
		params.ContentType = aws.String(input.ContentType)
	}
	out, err := m.client.rawClient.CreateMultipartUpload(ctx, params)
	if err != nil {
		return err
	}
	m.uploadID = out.UploadId
	return nil
}

func (m *multipartUpload) uploadPart(ctx context.Context, data []byte) error {
	if len(m.parts) >= maxMultipartParts {
		return errTooManyParts
	}
	number := aws.Int32(int32(len(m.parts) + 1))
	out, err := m.client.rawClient.UploadPart(ctx, &s3.UploadPartInput{
		Bucket: m.bucket, Key: m.key, UploadId: m.uploadID, PartNumber: number, Body: bytes.NewReader(data),
		ChecksumAlgorithm: m.checksumAlgorithm(),
	})
	if err != nil {
		return err
	}
	m.parts = append(m.parts, newCompletedPart(number, out))
	return nil
}

func (m *multipartUpload) uploadChunk(ctx context.Context, rest io.Reader) (bool, error) {
	chunk, err := readChunk(rest, streamPartSize)
	if err != nil || len(chunk) == 0 {
		return true, err
	}
	return false, m.uploadPart(ctx, chunk)
}

func (m *multipartUpload) uploadRemaining(ctx context.Context, rest io.Reader) error {
	done, err := m.uploadChunk(ctx, rest)
	if done || err != nil {
		return err
	}
	return m.uploadRemaining(ctx, rest)
}

func (m *multipartUpload) run(ctx context.Context, first []byte, rest io.Reader) (*contracts.PutObjectOutput, error) {
	if err := m.uploadPart(ctx, first); err != nil {
		return nil, err
	}
	if err := m.uploadRemaining(ctx, rest); err != nil {
		return nil, err
	}
	return m.complete(ctx)
}
