package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

type initProviderCase struct {
	configure    func(*config.Config)
	name         string
	provider     string
	bucket       string
	region       string
	endpoint     string
	checksumMode string
	usePathStyle bool
}

const defaultTestRegion = "us-east-1"

func initProviderTestCases() []initProviderCase {
	return []initProviderCase{
		{
			name: "aws", provider: "aws", bucket: "my-aws-bucket",
			region: defaultTestRegion, endpoint: "", usePathStyle: false, checksumMode: "when_supported",
			configure: func(cfg *config.Config) { cfg.Bucket = "my-aws-bucket"; cfg.Region = defaultTestRegion },
		},
		{
			name: "gcs", provider: "gcs", bucket: "my-gcs-bucket",
			region: "auto", endpoint: "https://storage.googleapis.com",
			usePathStyle: true, checksumMode: "when_required",
			configure: func(cfg *config.Config) {
				cfg.Bucket = "my-gcs-bucket"
				cfg.Region = "auto"
				cfg.Endpoint = "https://storage.googleapis.com"
			},
		},
		{
			name: "s3-compatible", provider: "s3-compatible", bucket: "my-s3-bucket",
			region: defaultTestRegion, endpoint: "https://minio.corp:9000",
			usePathStyle: true, checksumMode: "when_required",
			configure: func(cfg *config.Config) {
				cfg.Bucket = "my-s3-bucket"
				cfg.Region = defaultTestRegion
				cfg.Endpoint = "https://minio.corp:9000"
			},
		},
	}
}

func assertProviderInitFields(t *testing.T, initLog map[string]any, tc initProviderCase) {
	t.Helper()
	assert.Equal(t, tc.provider, initLog["provider"])
	assert.Equal(t, tc.bucket, initLog["bucket"])
	assert.Equal(t, tc.region, initLog["region"])
	assert.Equal(t, tc.endpoint, initLog["endpoint"])
	assert.Equal(t, tc.usePathStyle, initLog["use_path_style"])
	assert.Equal(t, tc.checksumMode, initLog["checksum_mode"])
}

func captureProviderInitLogs(t *testing.T, tc initProviderCase, mode storageLogMode) []map[string]any {
	t.Helper()
	isolateAWSEnv(t)
	var buf bytes.Buffer
	setStorageTestLogger(t, mode.handler(&buf))
	cfg := config.DefaultConfig()
	tc.configure(cfg)
	_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
	require.NoError(t, err)
	return parseStorageLogs(t, buf.Bytes())
}

func assertProviderInitLogs(t *testing.T, tc initProviderCase, mode storageLogMode) {
	t.Helper()
	entries := captureProviderInitLogs(t, tc, mode)
	initLog := findStorageLog(entries, "initializing storage client")
	require.NotNil(t, initLog)
	assertProviderInitFields(t, initLog, tc)
}

func TestConnect_Providers_InitLifecycleLogs(t *testing.T) {
	for _, mode := range storageLogModes() {
		for _, tc := range initProviderTestCases() {
			t.Run(mode.name+"_"+tc.name, func(t *testing.T) {
				assertProviderInitLogs(t, tc, mode)
			})
		}
	}
}
