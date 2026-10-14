// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import "github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"

// Public type aliases for object stream output and metadata retrieval.
type (
	// GetObjectOutput contains the downloaded object stream and metadata.
	GetObjectOutput = contracts.GetObjectOutput
	// HeadObjectInput specifies parameters for retrieving object metadata without body.
	HeadObjectInput = contracts.HeadObjectInput
	// HeadObjectOutput contains object metadata and attributes.
	HeadObjectOutput = contracts.HeadObjectOutput
)
