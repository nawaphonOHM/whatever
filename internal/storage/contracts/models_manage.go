// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import "time"

// DeleteObjectInput specifies parameters for deleting an object.
type DeleteObjectInput struct {
	Bucket string
	Key    string
}

// DeleteObjectOutput contains results from DeleteObject.
type DeleteObjectOutput struct {
	VersionID string
}

// ObjectSummary provides metadata overview for a listed object.
type ObjectSummary struct {
	LastModified time.Time
	Key          string
	ETag         string
	Size         int64
}
