// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import "github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"

// Public type aliases for object deletion and listing summaries.
type (
	// DeleteObjectInput specifies parameters for deleting an object.
	DeleteObjectInput = contracts.DeleteObjectInput
	// DeleteObjectOutput contains results from DeleteObject.
	DeleteObjectOutput = contracts.DeleteObjectOutput
	// ObjectSummary provides metadata overview for a listed object.
	ObjectSummary = contracts.ObjectSummary
)
