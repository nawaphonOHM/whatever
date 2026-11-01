package client

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

const (
	calcWhenRequired  = "WHEN_REQUIRED"
	calcWhenSupported = "WHEN_SUPPORTED"
)

func compatMatrixCasesGCSAndS3() []debugCompatCase {
	return []debugCompatCase{
		{
			name: "gcs_defaults",
			configure: func(cfg *config.Config) {
				cfg.Endpoint = "https://storage.googleapis.com"
			},
			expected: debugCompatExpected{
				provider: "gcs", reqChecksumCalc: calcWhenRequired,
				respChecksumVal: calcWhenRequired, checksumBehavior: BehaviorChecksumSuppressed,
				usePathStyle: true, pathStyleBehavior: BehaviorPathStyleEnforced,
				disableSkipLog: true, authMode: AuthModeDefaultChain,
				hasDriverOverrides: false,
			},
		},
		{
			name: "s3_compatible_static_creds",
			configure: func(cfg *config.Config) {
				cfg.Region = "us-east-1"
				cfg.Endpoint = "https://minio.corp:9000"
				cfg.AccessKeyID = "minio-key"
				cfg.SecretAccessKey = "minio-secret"
			},
			expected: debugCompatExpected{
				provider: "s3-compatible", reqChecksumCalc: calcWhenRequired,
				respChecksumVal: calcWhenRequired, checksumBehavior: BehaviorChecksumEnforced,
				usePathStyle: true, pathStyleBehavior: BehaviorPathStyleEnforced,
				disableSkipLog: true, authMode: AuthModeStaticCredentials,
				hasDriverOverrides: false,
			},
		},
	}
}

func compatMatrixCasesAWS() []debugCompatCase {
	return []debugCompatCase{
		{
			name: "aws_defaults",
			configure: func(cfg *config.Config) {
				cfg.Region = "us-east-1"
			},
			expected: debugCompatExpected{
				provider: "aws", reqChecksumCalc: calcWhenSupported,
				respChecksumVal: calcWhenSupported, checksumBehavior: BehaviorChecksumStandard,
				usePathStyle: false, pathStyleBehavior: BehaviorPathStyleStandard,
				disableSkipLog: true, authMode: AuthModeDefaultChain,
				hasDriverOverrides: false,
			},
		},
		{
			name: "aws_configured_path_style_and_driver_opt",
			configure: func(cfg *config.Config) {
				cfg.Region = "us-east-1"
				cfg.UsePathStyle = true
			},
			opts: []Option{WithDriverOptions(func(*s3.Options) {})},
			expected: debugCompatExpected{
				provider: "aws", reqChecksumCalc: calcWhenSupported,
				respChecksumVal: calcWhenSupported, checksumBehavior: BehaviorChecksumStandard,
				usePathStyle: true, pathStyleBehavior: BehaviorPathStyleConfigured,
				disableSkipLog: true, authMode: AuthModeDefaultChain,
				hasDriverOverrides: true,
			},
		},
	}
}
