package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

const (
	secSecretKey    = "SUPER_SECRET_KEY_999"
	secSessionToken = "SESSION_TOKEN_SENSITIVE_SENTINEL_777"
	secToken        = "token_xyz_do_not_log"
	secPassword     = "pass_sensitive_123"
	secFragment     = "frag_secret_part"
)

func securityTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Region = "us-east-1"
	cfg.AccessKeyID = "test-access-key-id"
	cfg.SecretAccessKey = secSecretKey
	cfg.SessionToken = secSessionToken
	cfg.Endpoint = "https://admin:" + secPassword + "@minio.internal:9000/data?" + secToken + "#" + secFragment
	cfg.BucketURL = "https://user:" + secPassword + "@bucket.storage.googleapis.com?" + secToken + "#" + secFragment
	return cfg
}

func assertNoSensitiveDataInOutput(t *testing.T, rawOutput string) {
	t.Helper()
	assert.NotContains(t, rawOutput, secSecretKey)
	assert.NotContains(t, rawOutput, secSessionToken)
	assert.NotContains(t, rawOutput, secToken)
	assert.NotContains(t, rawOutput, secPassword)
	assert.NotContains(t, rawOutput, secFragment)
}

func assertSecurityLogEntries(t *testing.T, entries []map[string]any) {
	t.Helper()
	detectLog := findStorageLog(entries, "evaluating storage provider detection")
	require.NotNil(t, detectLog)
	assert.Equal(t, "https://minio.internal:9000/data", detectLog["raw_endpoint"])
	assert.Equal(t, "https://bucket.storage.googleapis.com", detectLog["raw_bucket_url"])
	assert.Equal(t, "minio.internal", detectLog["inspected_host"])
	assert.Equal(t, "s3-compatible", detectLog["detected_provider"])

	compatLog := findStorageLog(entries, "storage client compatibility configuration")
	require.NotNil(t, compatLog)
	assert.Equal(t, AuthModeStaticCredentials, compatLog["auth_mode"])
}

func TestConnect_DebugAndInfo_NeverLogsCredentials(t *testing.T) {
	for _, mode := range storageLogModes() {
		t.Run(mode.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, mode.debugHandler(&buf))
			_, err := ConnectWithConfig(context.Background(), securityTestConfig(), WithSkipProbe(true))
			require.NoError(t, err)

			assertNoSensitiveDataInOutput(t, buf.String())
			assertSecurityLogEntries(t, parseStorageLogs(t, buf.Bytes()))
		})
	}
}
