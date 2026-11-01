package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	gcsPSCBase      = "https://storage-psc.p.googleapis.com"
	gcsCustomBase   = "https://storage-custom.p.googleapis.com"
	gcsRegionalBase = "https://storage.us-central1.rep.googleapis.com"
	awsExpressBase  = "https://s3express-use1-az4.us-east-1.amazonaws.com"
	awsVPCEBase     = "https://bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com"
	otherBucket     = "other"
)

// endpointCfg and bucketURLCfg build a config that targets host through Endpoint or BucketURL.
func endpointCfg(host, bucket string) *config.Config {
	return &config.Config{Bucket: bucket, Endpoint: "https://" + host}
}

func bucketURLCfg(host, bucket string) *config.Config {
	return &config.Config{Bucket: bucket, BucketURL: "https://" + host}
}

type routingCase struct {
	host      string
	wantBase  string
	bucket    string
	pathStyle bool
}

type mismatchCase struct {
	host       string
	wantBucket string
}

func bucketOrDefault(bucket string) string {
	if bucket == "" {
		return testBucket
	}
	return bucket
}

func assertBaseEndpoint(t *testing.T, got *string, want string) {
	t.Helper()
	if want == "" {
		assert.Nil(t, got)
		return
	}
	require.NotNil(t, got)
	assert.Equal(t, want, *got)
}

func runRoutingCases(
	t *testing.T, provider contracts.Provider, mk func(host, bucket string) *config.Config, cases []routingCase,
) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.host, func(t *testing.T) {
			t.Parallel()
			opts, err := client.BuildS3Options(mk(tt.host, bucketOrDefault(tt.bucket)), provider, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.pathStyle, opts.UsePathStyle)
			assertBaseEndpoint(t, opts.BaseEndpoint, tt.wantBase)
		})
	}
}

func runMismatchCases(
	t *testing.T, provider contracts.Provider, mk func(host, bucket string) *config.Config, cases []mismatchCase,
) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.host, func(t *testing.T) {
			t.Parallel()
			opts, err := client.BuildS3Options(mk(tt.host, testBucket), provider, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(),
				`bucket "mybucket" in configuration does not match bucket "`+tt.wantBucket+`" in URL`)
			assert.Empty(t, opts)
		})
	}
}
