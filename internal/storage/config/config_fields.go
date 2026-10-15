package config

import (
	"time"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// Port boundary constants for storage endpoint.
const (
	MinPort = 1
	MaxPort = 65535
)

// Bucket name length boundary constants (3 to 222 covers AWS S3 and GCS rules).
const (
	MinBucketLength = 3
	MaxBucketLength = 222
)

// MaxPresignExpiry is the maximum allowed duration for presigned URLs (SigV4 limit: 7 days).
const (
	maxPresignExpiryDays = 7
	hoursPerDay          = 24
	MaxPresignExpiry     = maxPresignExpiryDays * hoursPerDay * time.Hour
)

// ErrNilConfig is returned when validating a nil configuration.
var ErrNilConfig = contracts.ErrNilConfig

// BaseFields defines core storage configuration including bucket, region, and credentials.
type BaseFields struct {
	Bucket          string `env:"OHM9996_STORAGE_BUCKET"`
	Region          string `env:"OHM9996_STORAGE_REGION" envDefault:"us-east-1"`
	AccessKeyID     string `env:"OHM9996_STORAGE_ACCESS_KEY_ID" envDefault:""`
	SecretAccessKey string `env:"OHM9996_STORAGE_SECRET_ACCESS_KEY" envDefault:""`
	SessionToken    string `env:"OHM9996_STORAGE_SESSION_TOKEN" envDefault:""`
}

// Config defines the complete configuration for object storage operations.
type Config struct {
	BaseFields
	EndpointFields
	TimeoutFields
}
