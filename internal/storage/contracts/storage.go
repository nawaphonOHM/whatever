// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Provider represents the cloud object storage provider type.
type Provider string

const (
	// ProviderAWS represents native Amazon Web Services S3.
	ProviderAWS Provider = "aws"
	// ProviderGCS represents Google Cloud Storage XML API (S3-compatible).
	ProviderGCS Provider = "gcs"
	// ProviderS3Compatible represents generic S3-compatible storage providers (e.g. MinIO, LocalStack).
	ProviderS3Compatible Provider = "s3-compatible"
)

// String returns the string representation of the Provider.
func (p Provider) String() string {
	return string(p)
}

// IsValid checks if the provider is one of the supported provider types.
func (p Provider) IsValid() bool {
	switch p {
	case ProviderAWS, ProviderGCS, ProviderS3Compatible:
		return true
	default:
		return false
	}
}

// Sentinel errors for cloud object storage operations.
var (
	// ErrObjectNotFound is returned when the requested object key does not exist.
	ErrObjectNotFound = errors.New("storage object not found")
	// ErrBucketNotFound is returned when the target bucket does not exist.
	ErrBucketNotFound = errors.New("storage bucket not found")
	// ErrNilClient is returned when an operation is performed on an uninitialized client.
	ErrNilClient = errors.New("storage client is not initialized")
	// ErrNilConfig is returned when operations receive a nil configuration.
	ErrNilConfig = errors.New("storage config cannot be nil")
	// ErrAccessDenied is returned when the storage service rejects the request (HTTP 403).
	ErrAccessDenied = errors.New("storage access denied")
	// ErrInvalidExpiry is returned when a pre-signed URL expiration is negative or exceeds 7 days.
	ErrInvalidExpiry = errors.New("storage presign expiry is invalid")
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
