// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import "io"

// PutObjectInput specifies parameters for uploading an object.
type PutObjectInput struct {
	Body        io.Reader
	Metadata    map[string]string
	Bucket      string
	Key         string
	ContentType string
}

// PutObjectOutput contains results from PutObject.
type PutObjectOutput struct {
	ETag      string
	VersionID string
}

// GetObjectInput specifies parameters for downloading an object.
type GetObjectInput struct {
	Bucket string
	Key    string
}
