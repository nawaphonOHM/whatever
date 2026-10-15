// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import "github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"

// Public type aliases for presigned URL generation and results.
type (
	// PresignGetObjectInput specifies parameters for generating pre-signed GET URL.
	PresignGetObjectInput = contracts.PresignGetObjectInput
	// PresignPutObjectInput specifies parameters for generating pre-signed PUT URL.
	PresignPutObjectInput = contracts.PresignPutObjectInput
	// PresignedURL contains the pre-signed URL and request metadata.
	PresignedURL = contracts.PresignedURL
)
