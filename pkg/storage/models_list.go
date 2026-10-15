// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import "github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"

// Public type aliases for bucket listing parameters and results.
type (
	// ListObjectsInput specifies parameters for listing objects in a bucket.
	ListObjectsInput = contracts.ListObjectsInput
	// ListObjectsOutput contains listed objects and common prefixes.
	ListObjectsOutput = contracts.ListObjectsOutput
)
