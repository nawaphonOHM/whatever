package client

import (
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
	secHashFrag     = "#" + secFragment
	secAt           = "@"

	secTestBucket    = "test-bucket"
	secTestRegion    = "us-east-1"
	secTestAccessKey = "test-access-key-id"

	secHostMinio   = "minio.internal"
	secAdminPrefix = "https://admin:"
	secUserPrefix  = "user:"
	secProviderS3  = "s3-compatible"

	secLogEvalDetect  = "evaluating storage provider detection"
	secLogInspectHost = "inspected_host"
	secLogDetectProv  = "detected_provider"
)

func securityTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = secTestBucket
	cfg.Region = secTestRegion
	cfg.AccessKeyID = secTestAccessKey
	cfg.SecretAccessKey = secSecretKey
	cfg.SessionToken = secSessionToken
	cfg.Endpoint = secAdminPrefix + secPassword + secAt + secHostMinio + ":9000/data?" + secToken + secHashFrag
	cfg.BucketURL = "https://user:" + secPassword + "@bucket.storage.googleapis.com?" + secToken + secHashFrag
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

func assertErrorNoSecrets(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.NotContains(t, err.Error(), secSecretKey)
	assert.NotContains(t, err.Error(), secPassword)
	assert.NotContains(t, err.Error(), secToken)
}

func assertMalformedRejection(t *testing.T, cfg *config.Config, rawOutput string, connErr error) {
	t.Helper()
	assertErrorNoSecrets(t, connErr)
	assertErrorNoSecrets(t, cfg.Validate())
	assertNoSensitiveDataInOutput(t, rawOutput)
}

func assertSecurityLogEntries(t *testing.T, entries []map[string]any) {
	t.Helper()
	detectLog := findStorageLog(entries, secLogEvalDetect)
	require.NotNil(t, detectLog)
	assert.Equal(t, "https://minio.internal:9000/data", detectLog["raw_endpoint"])
	assert.Equal(t, "https://bucket.storage.googleapis.com", detectLog["raw_bucket_url"])
	assert.Equal(t, secHostMinio, detectLog[secLogInspectHost])
	assert.Equal(t, secProviderS3, detectLog[secLogDetectProv])

	compatLog := findStorageLog(entries, "storage client compatibility configuration")
	require.NotNil(t, compatLog)
	assert.Equal(t, AuthModeStaticCredentials, compatLog["auth_mode"])
}
