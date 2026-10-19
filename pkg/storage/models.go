// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import "github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"

// Public type aliases for object storage operations, inputs, outputs, and metadata summaries.
type (
	// PutObjectInput specifies parameters for uploading an object.
	PutObjectInput = contracts.PutObjectInput
	// PutObjectOutput contains results from PutObject.
	PutObjectOutput = contracts.PutObjectOutput
	// GetObjectInput specifies parameters for downloading an object.
	GetObjectInput = contracts.GetObjectInput
	// GetObjectOutput contains the downloaded object stream and metadata.
	GetObjectOutput = contracts.GetObjectOutput
	// HeadObjectInput specifies parameters for retrieving object metadata without body.
	HeadObjectInput = contracts.HeadObjectInput
	// HeadObjectOutput contains object metadata and attributes.
	HeadObjectOutput = contracts.HeadObjectOutput
	// DeleteObjectInput specifies parameters for deleting an object.
	DeleteObjectInput = contracts.DeleteObjectInput
	// DeleteObjectOutput contains results from DeleteObject.
	DeleteObjectOutput = contracts.DeleteObjectOutput
	// ListObjectsInput specifies parameters for listing objects in a bucket.
	ListObjectsInput = contracts.ListObjectsInput
	// ListObjectsOutput contains listed objects and common prefixes.
	ListObjectsOutput = contracts.ListObjectsOutput
	// ObjectSummary provides metadata overview for a listed object.
	ObjectSummary = contracts.ObjectSummary
	// PresignGetObjectInput specifies parameters for generating pre-signed GET URL.
	PresignGetObjectInput = contracts.PresignGetObjectInput
	// PresignPutObjectInput specifies parameters for generating pre-signed PUT URL.
	PresignPutObjectInput = contracts.PresignPutObjectInput
	// PresignedURL contains the pre-signed URL and request metadata.
	PresignedURL = contracts.PresignedURL
)
