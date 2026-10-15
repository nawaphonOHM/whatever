// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import (
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// Storage defines the abstract contract interface for cloud object storage operations.
type Storage = contracts.Storage

// Provider represents the cloud object storage provider type.
type Provider = contracts.Provider

// Provider constants for supported cloud object storage systems.
const (
	// ProviderAWS represents native Amazon Web Services S3.
	ProviderAWS = contracts.ProviderAWS
	// ProviderGCS represents Google Cloud Storage XML API (S3-compatible).
	ProviderGCS = contracts.ProviderGCS
	// ProviderS3Compatible represents generic S3-compatible storage providers (e.g. MinIO, LocalStack).
	ProviderS3Compatible = contracts.ProviderS3Compatible
)

// Sentinel errors for cloud storage operations.
var (
	// ErrObjectNotFound is returned when the requested object key does not exist.
	ErrObjectNotFound = contracts.ErrObjectNotFound
	// ErrBucketNotFound is returned when the target bucket does not exist.
	ErrBucketNotFound = contracts.ErrBucketNotFound
	// ErrNilClient is returned when an operation is performed on an uninitialized client.
	ErrNilClient = contracts.ErrNilClient
	// ErrNilConfig is returned when operations receive a nil configuration.
	ErrNilConfig = contracts.ErrNilConfig
	// ErrAccessDenied is returned when the storage service rejects the request (HTTP 403).
	ErrAccessDenied = contracts.ErrAccessDenied
	// ErrInvalidExpiry is returned when a pre-signed URL expiration is negative or exceeds 7 days.
	ErrInvalidExpiry = contracts.ErrInvalidExpiry
)
