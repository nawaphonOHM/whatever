// Package contracts defines abstract interfaces, data contracts, and errors
// for cloud object storage operations.
package contracts

// Provider represents the cloud object storage provider type.
type Provider string

const (
	// ProviderAWS represents native Amazon Web Services S3.
	ProviderAWS Provider = "aws"
	// ProviderGCS represents Google Cloud Storage XML API (S3-compatible).
	ProviderGCS Provider = "gcs"
	// ProviderS3Compatible represents generic S3-compatible storage providers (e.g. MinIO, LocalStack).
	ProviderS3Compatible Provider = "s3-compatible"
)

// String returns the string representation of the Provider.
func (p Provider) String() string {
	return string(p)
}

// IsValid checks if the provider is one of the supported provider types.
func (p Provider) IsValid() bool {
	switch p {
	case ProviderAWS, ProviderGCS, ProviderS3Compatible:
		return true
	default:
		return false
	}
}
