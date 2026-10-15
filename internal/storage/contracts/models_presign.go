// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

import "time"

// PresignGetObjectInput specifies parameters for generating pre-signed GET URL.
type PresignGetObjectInput struct {
	Bucket    string
	Key       string
	ExpiresIn time.Duration
}

// PresignPutObjectInput specifies parameters for generating pre-signed PUT URL.
type PresignPutObjectInput struct {
	Bucket      string
	Key         string
	ContentType string
	ExpiresIn   time.Duration
}

// PresignedURL contains the pre-signed URL and request metadata.
type PresignedURL struct {
	Headers   map[string][]string
	URL       string
	Method    string
	ExpiresIn time.Duration
}
