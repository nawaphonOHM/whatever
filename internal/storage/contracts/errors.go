// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import "errors"

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
