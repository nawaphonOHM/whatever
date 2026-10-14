// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Storage defines the abstract contract interface for cloud object storage operations.
type Storage interface {
	// PutObject uploads an object from an io.Reader stream with optional metadata.
	// Any io.Reader is accepted: seekable bodies are sent as-is, other streams are buffered in
	// bounded parts (multipart upload for large or unknown-length streams).
	PutObject(context.Context, PutObjectInput) (*PutObjectOutput, error)
	// GetObject retrieves an object stream and its associated metadata.
	// The caller owns GetObjectOutput.Body and must Close it. The request timeout covers the time
	// until response headers arrive; reading the body is bounded only by the caller's context.
	GetObject(context.Context, GetObjectInput) (*GetObjectOutput, error)
	// DeleteObject deletes an object from the specified bucket.
	DeleteObject(context.Context, DeleteObjectInput) (*DeleteObjectOutput, error)
	// HeadObject retrieves object metadata without downloading the payload body.
	HeadObject(context.Context, HeadObjectInput) (*HeadObjectOutput, error)
	// ListObjects enumerates object summaries and common prefixes in a bucket.
	ListObjects(context.Context, ListObjectsInput) (*ListObjectsOutput, error)
	// Exists checks whether a specific object key exists in the bucket.
	// A missing object yields (false, nil); a missing bucket or access denial yields an error.
	Exists(context.Context, string, string) (bool, error)
	// PresignGetObject generates a pre-signed URL for downloading an object.
	// ExpiresIn of 0 uses the configured default; values above 7 days are rejected (ErrInvalidExpiry).
	PresignGetObject(context.Context, PresignGetObjectInput) (*PresignedURL, error)
	// PresignPutObject generates a pre-signed URL for uploading an object.
	PresignPutObject(context.Context, PresignPutObjectInput) (*PresignedURL, error)
	// Ping probes connectivity and accessibility of the object storage service.
	Ping(context.Context) error
	// Provider returns the cloud provider type (AWS, GCS, or S3-compatible).
	Provider() Provider
	// DefaultBucket returns the default bucket configured for this storage client.
	DefaultBucket() string
	// RawClient returns the underlying AWS SDK v2 S3 client handle.
	RawClient() *s3.Client
}
