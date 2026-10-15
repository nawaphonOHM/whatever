// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func newCompletedPart(number *int32, out *s3.UploadPartOutput) types.CompletedPart {
	return types.CompletedPart{
		ChecksumCRC32:     out.ChecksumCRC32,
		ChecksumCRC32C:    out.ChecksumCRC32C,
		ChecksumCRC64NVME: out.ChecksumCRC64NVME,
		ChecksumSHA1:      out.ChecksumSHA1,
		ChecksumSHA256:    out.ChecksumSHA256,
		ETag:              out.ETag,
		PartNumber:        number,
	}
}

func (m *multipartUpload) complete(ctx context.Context) (*contracts.PutObjectOutput, error) {
	out, err := m.client.rawClient.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket: m.bucket, Key: m.key, UploadId: m.uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{Parts: m.parts},
	})
	if err != nil {
		return nil, err
	}
	return &contracts.PutObjectOutput{ETag: aws.ToString(out.ETag), VersionID: aws.ToString(out.VersionId)}, nil
}

// abort discards uploaded parts, using a context detached from the failed request's cancellation.
func (m *multipartUpload) abort(ctx context.Context) error {
	abortCtx, cancel := withTimeout(context.WithoutCancel(ctx), m.client.requestTimeout)
	defer cancel()
	_, err := m.client.rawClient.AbortMultipartUpload(abortCtx, &s3.AbortMultipartUploadInput{
		Bucket: m.bucket, Key: m.key, UploadId: m.uploadID,
	})
	return err
}

// putMultipart uploads first and the remaining stream as parts, aborting the upload on any failure.
func (c *Client) putMultipart(
	ctx context.Context,
	input contracts.PutObjectInput,
	first []byte,
	rest io.Reader,
) (*contracts.PutObjectOutput, error) {
	m := &multipartUpload{
		client: c, bucket: aws.String(c.resolveBucket(input.Bucket)), key: aws.String(input.Key),
	}
	if err := m.create(ctx, input); err != nil {
		return nil, err
	}
	out, err := m.run(ctx, first, rest)
	if err != nil {
		return nil, errors.Join(err, m.abort(ctx))
	}
	return out, nil
}
