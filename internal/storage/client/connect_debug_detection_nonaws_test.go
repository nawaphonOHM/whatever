package client

import (
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

const gcsGlobalEndpoint = "https://storage.googleapis.com"

func detectionCasesNonAWS() []debugDetectionCase {
	return []debugDetectionCase{
		{
			name: "gcs_hostname",
			configure: func(cfg *config.Config) {
				cfg.Endpoint = gcsGlobalEndpoint
			},
			expected: debugDetectionExpected{
				rawEndpoint: gcsGlobalEndpoint, rawBucketURL: "",
				inspectedHost: "storage.googleapis.com", detectedProvider: "gcs",
				detectionRule: detector.RuleGCSHostname, extractedBucket: "",
				normalizedEndpoint: gcsGlobalEndpoint,
			},
		},
		{
			name: "gcs_virtual_hosted",
			configure: func(cfg *config.Config) {
				cfg.Bucket = "my-gcs-bucket"
				cfg.BucketURL = "https://my-gcs-bucket.storage.googleapis.com"
			},
			expected: debugDetectionExpected{
				rawEndpoint: "", rawBucketURL: "https://my-gcs-bucket.storage.googleapis.com",
				inspectedHost: "my-gcs-bucket.storage.googleapis.com", detectedProvider: "gcs",
				detectionRule: detector.RuleGCSHostname, extractedBucket: "my-gcs-bucket",
				normalizedEndpoint: gcsGlobalEndpoint,
			},
		},
		{
			name: "s3_compatible_fallback",
			configure: func(cfg *config.Config) {
				cfg.Region = "us-east-1"
				cfg.Endpoint = "https://custom.minio.corp:9000"
			},
			expected: debugDetectionExpected{
				rawEndpoint: "https://custom.minio.corp:9000", rawBucketURL: "",
				inspectedHost: "custom.minio.corp", detectedProvider: "s3-compatible",
				detectionRule: detector.RuleS3CompatibleFallback, extractedBucket: "",
				normalizedEndpoint: "https://custom.minio.corp:9000",
			},
		},
	}
}
