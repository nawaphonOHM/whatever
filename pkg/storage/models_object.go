// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import "github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"

// Public type aliases for object upload and retrieval inputs and outputs.
type (
	// PutObjectInput specifies parameters for uploading an object.
	PutObjectInput = contracts.PutObjectInput
	// PutObjectOutput contains results from PutObject.
	PutObjectOutput = contracts.PutObjectOutput
	// GetObjectInput specifies parameters for downloading an object.
	GetObjectInput = contracts.GetObjectInput
)
