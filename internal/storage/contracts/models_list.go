// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

// ListObjectsInput specifies parameters for listing objects in a bucket.
type ListObjectsInput struct {
	Bucket    string
	Prefix    string
	Delimiter string
	MaxKeys   int32
}

// ListObjectsOutput contains listed objects and common prefixes.
type ListObjectsOutput struct {
	Objects  []ObjectSummary
	Prefixes []string
}
