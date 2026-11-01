package client_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestBuildS3Options_AWSEndpointRouting(t *testing.T) {
	t.Parallel()
	runRoutingCases(t, contracts.ProviderAWS, endpointCfg, []routingCase{
		{
			"MYBUCKET.Bucket.VPCE-1a2b-3c4d.S3.US-EAST-1.VPCE.AMAZONAWS.COM",
			"https://Bucket.VPCE-1a2b-3c4d.S3.US-EAST-1.VPCE.AMAZONAWS.COM", "", false,
		},
		{"mybucket.s3express-use1-az4.us-east-1.amazonaws.com", awsExpressBase, "", false},
		{"s3express-use1-az4.us-east-1.amazonaws.com", awsExpressBase, "", false},
		{"dynamodb.us-east-1.amazonaws.com", "https://dynamodb.us-east-1.amazonaws.com", "", true},
		{"vpce-logs.s3.us-east-1.amazonaws.com", "", "vpce-logs", false},
		{"bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", awsVPCEBase, "", false},
	})
}

func TestBuildS3Options_AWSBucketMismatch(t *testing.T) {
	t.Parallel()
	runMismatchCases(t, contracts.ProviderAWS, endpointCfg, []mismatchCase{
		{"s3express-x.s3.us-east-1.amazonaws.com", "s3express-x"},
		{"vpce-logs.s3.us-east-1.amazonaws.com", "vpce-logs"},
		{"other.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", otherBucket},
		{"other.s3.us-east-1.amazonaws.com.", otherBucket},
	})
}

func TestBuildS3Options_AWSBucketURLMismatch(t *testing.T) {
	t.Parallel()
	runMismatchCases(t, contracts.ProviderAWS, bucketURLCfg, []mismatchCase{
		{"other.s3.us-east-1.amazonaws.com.", otherBucket},
	})
}

func TestBuildS3Options_AWSBucketURLMixedCaseVPCE(t *testing.T) {
	t.Parallel()
	runRoutingCases(t, contracts.ProviderAWS, bucketURLCfg, []routingCase{
		{
			"MYBUCKET.Bucket.VPCE-1a2b-3c4d.S3.US-EAST-1.VPCE.AMAZONAWS.COM",
			"https://Bucket.VPCE-1a2b-3c4d.S3.US-EAST-1.VPCE.AMAZONAWS.COM", "", false,
		},
	})
}
