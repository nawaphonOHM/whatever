package client

import (
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

const providerAWSName = "aws"

func detectionCasesAWS() []debugDetectionCase {
	return []debugDetectionCase{
		{
			name: "default_aws",
			configure: func(cfg *config.Config) {
				cfg.Region = "us-east-1"
			},
			expected: debugDetectionExpected{
				rawEndpoint: "", rawBucketURL: "", inspectedHost: "",
				detectedProvider: providerAWSName, detectionRule: detector.RuleDefaultAWS,
				extractedBucket: "", normalizedEndpoint: "",
			},
		},
		{
			name: "aws_hostname",
			configure: func(cfg *config.Config) {
				cfg.Region = "us-west-2"
				cfg.Endpoint = "https://s3.us-west-2.amazonaws.com"
			},
			expected: debugDetectionExpected{
				rawEndpoint: "https://s3.us-west-2.amazonaws.com", rawBucketURL: "",
				inspectedHost: "s3.us-west-2.amazonaws.com", detectedProvider: providerAWSName,
				detectionRule: detector.RuleAWSHostname, extractedBucket: "",
				normalizedEndpoint: "",
			},
		},
		{
			name: "aws_virtual_hosted",
			configure: func(cfg *config.Config) {
				cfg.Bucket = "my-bucket"
				cfg.BucketURL = "https://my-bucket.s3.amazonaws.com"
			},
			expected: debugDetectionExpected{
				rawEndpoint: "", rawBucketURL: "https://my-bucket.s3.amazonaws.com",
				inspectedHost: "my-bucket.s3.amazonaws.com", detectedProvider: providerAWSName,
				detectionRule: detector.RuleAWSHostname, extractedBucket: "my-bucket",
				normalizedEndpoint: "",
			},
		},
		{
			name: "aws_vpce_custom",
			configure: func(cfg *config.Config) {
				cfg.Bucket = "mybucket"
				cfg.Region = "us-east-1"
				cfg.Endpoint = "https://mybucket.vpce-1a.s3.us-east-1.vpce.amazonaws.com"
			},
			expected: debugDetectionExpected{
				rawEndpoint: "https://mybucket.vpce-1a.s3.us-east-1.vpce.amazonaws.com", rawBucketURL: "",
				inspectedHost: "mybucket.vpce-1a.s3.us-east-1.vpce.amazonaws.com", detectedProvider: providerAWSName,
				detectionRule: detector.RuleAWSHostname, extractedBucket: "mybucket",
				normalizedEndpoint: "https://vpce-1a.s3.us-east-1.vpce.amazonaws.com",
			},
		},
	}
}
