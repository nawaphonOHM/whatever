// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import (
	"io"
	"time"
)

// GetObjectOutput contains the downloaded object stream and metadata.
type GetObjectOutput struct {
	LastModified  time.Time
	Body          io.ReadCloser
	Metadata      map[string]string
	ContentType   string
	ETag          string
	ContentLength int64
}

// HeadObjectInput specifies parameters for retrieving object metadata without body.
type HeadObjectInput struct {
	Bucket string
	Key    string
}

// HeadObjectOutput contains object metadata and attributes.
type HeadObjectOutput struct {
	LastModified  time.Time
	Metadata      map[string]string
	ContentType   string
	ETag          string
	ContentLength int64
}
