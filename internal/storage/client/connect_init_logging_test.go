package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	testLogBucket   = "test-bucket"
	testLogRegion   = "ap-southeast-2"
	testLogEndpoint = "https://storage.example.test"
)

func TestStorageInitializationAndCompletionLogs_JSONAndGCP(t *testing.T) {
	for _, mode := range storageLogModes() {
		t.Run(mode.name, func(t *testing.T) {
			assertStorageLifecycleLogs(t, mode)
		})
	}
}

func testLifecycleAWSParams() lifecycleLogParams {
	return lifecycleLogParams{
		provider: contracts.ProviderAWS,
		bucket:   testLogBucket,
		s3Opts: s3.Options{
			Region:       testLogRegion,
			BaseEndpoint: aws.String(testLogEndpoint),
			UsePathStyle: true,
		},
	}
}

func assertStorageLifecycleLogs(t *testing.T, mode storageLogMode) {
	t.Helper()
	isolateAWSEnv(t)
	var output bytes.Buffer
	setStorageTestLogger(t, mode.handler(&output))
	p := testLifecycleAWSParams()
	logInitialization(context.Background(), p)
	logStartupComplete(context.Background(), p)
	entries := parseStorageLogs(t, output.Bytes())
	assertStorageInitializationEntry(t, findStorageLog(entries, "initializing storage client"), mode)
	assertStorageCompletionEntry(t, findStorageLog(entries, "storage client initialization complete"), mode)
}

func assertStorageInitializationEntry(t *testing.T, entry map[string]any, mode storageLogMode) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, "INFO", entry[mode.levelKey])
	assert.Equal(t, "aws", entry["provider"])
	assert.Equal(t, testLogBucket, entry["bucket"])
	assert.Equal(t, testLogRegion, entry["region"])
	assert.Equal(t, testLogEndpoint, entry["endpoint"])
	assert.Equal(t, true, entry["use_path_style"])
	assert.Equal(t, "when_supported", entry["checksum_mode"])
}

func assertStorageCompletionEntry(t *testing.T, entry map[string]any, mode storageLogMode) {
	t.Helper()
	assertStorageInitializationEntry(t, entry, mode)
}
