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

func TestConnect_SchemelessWithQueryScheme_NeverLogsCredentials(t *testing.T) {
	for _, mode := range storageLogModes() {
		t.Run(mode.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, mode.debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = "test-bucket"
			cfg.Region = "us-east-1"
			cfg.AccessKeyID = "test-access-key-id"
			cfg.SecretAccessKey = secSecretKey
			cfg.SessionToken = secSessionToken
			cfg.Endpoint = "admin:" + secPassword + "@minio.internal:9000/data?redir=https://other.internal/x&token=" + secToken + "#" + secFragment
			cfg.BucketURL = "user:" + secPassword + "@custom-host:9000/data?r=a://b&token=" + secToken + "#frag://b"

			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.NoError(t, err)

			rawOutput := buf.String()
			assertNoSensitiveDataInOutput(t, rawOutput)

			entries := parseStorageLogs(t, buf.Bytes())
			detectLog := findStorageLog(entries, "evaluating storage provider detection")
			require.NotNil(t, detectLog)
			assert.Equal(t, "minio.internal:9000/data", detectLog["raw_endpoint"])
			assert.Equal(t, "custom-host:9000/data", detectLog["raw_bucket_url"])
			assert.Equal(t, "minio.internal", detectLog["inspected_host"])
			assert.Equal(t, "s3-compatible", detectLog["detected_provider"])
		})
	}
}

func TestConnect_AcceptedInputsWithSecretPayload_DiagnosticLogsSanitized(t *testing.T) {
	testCases := []struct {
		name              string
		endpoint          string
		bucketURL         string
		bucketName        string
		expectRawEndpoint string
		expectRawBucket   string
		expectHost        string
		expectProvider    string
	}{
		{
			name:              "absolute endpoint and bucket URL with credentials and query",
			endpoint:          "https://admin:" + secPassword + "@minio.internal:9000/data?token=" + secToken + "#" + secFragment,
			bucketURL:         "https://user:" + secPassword + "@bucket.storage.googleapis.com?token=" + secToken + "#" + secFragment,
			expectRawEndpoint: "https://minio.internal:9000/data",
			expectRawBucket:   "https://bucket.storage.googleapis.com",
			expectHost:        "minio.internal",
			expectProvider:    "s3-compatible",
		},
		{
			name:              "schemeless endpoint with query scheme separator and nested redirect",
			endpoint:          "admin:" + secPassword + "@minio.internal:9000/data?redir=https://other.internal/x&token=" + secToken + "#" + secFragment,
			bucketURL:         "user:" + secPassword + "@custom-host:9000/data?r=a://b&token=" + secToken + "#frag://b",
			expectRawEndpoint: "minio.internal:9000/data",
			expectRawBucket:   "custom-host:9000/data",
			expectHost:        "minio.internal",
			expectProvider:    "s3-compatible",
		},
		{
			name:              "schemeless aws bucket URL with credentials",
			endpoint:          "",
			bucketURL:         "user:" + secPassword + "@mybucket.s3.us-east-1.amazonaws.com/data?token=" + secToken,
			bucketName:        "mybucket",
			expectRawEndpoint: "",
			expectRawBucket:   "mybucket.s3.us-east-1.amazonaws.com/data",
			expectHost:        "mybucket.s3.us-east-1.amazonaws.com",
			expectProvider:    "aws",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, storageLogModes()[0].debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = "test-bucket"
			if tc.bucketName != "" {
				cfg.Bucket = tc.bucketName
			}
			cfg.Region = "us-east-1"
			cfg.AccessKeyID = "test-access-key-id"
			cfg.SecretAccessKey = secSecretKey
			cfg.SessionToken = secSessionToken
			cfg.Endpoint = tc.endpoint
			cfg.BucketURL = tc.bucketURL

			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.NoError(t, err)

			rawOutput := buf.String()
			assertNoSensitiveDataInOutput(t, rawOutput)

			entries := parseStorageLogs(t, buf.Bytes())
			detectLog := findStorageLog(entries, "evaluating storage provider detection")
			require.NotNil(t, detectLog, "diagnostic detection log must be emitted for accepted inputs")
			assert.Equal(t, tc.expectRawEndpoint, detectLog["raw_endpoint"])
			assert.Equal(t, tc.expectRawBucket, detectLog["raw_bucket_url"])
			assert.Equal(t, tc.expectHost, detectLog["inspected_host"])
			assert.Equal(t, tc.expectProvider, detectLog["detected_provider"])
		})
	}
}

func TestConnect_MalformedSecurityInputs_RejectionsDoNotLeakSecrets(t *testing.T) {
	malformedCases := []struct {
		name      string
		endpoint  string
		bucketURL string
	}{
		{
			name:      "whitespace in endpoint",
			endpoint:  "https://admin:" + secPassword + "@minio .internal:9000/data?" + secToken,
			bucketURL: "https://user:" + secPassword + "@storage.googleapis.com",
		},
		{
			name:      "tab in bucket URL",
			endpoint:  "",
			bucketURL: "user:" + secPassword + "@storage.googleapis\t.com?" + secToken,
		},
		{
			name:      "schemeless with only credentials",
			endpoint:  "user:" + secPassword + "@",
			bucketURL: "https://s3.amazonaws.com",
		},
		{
			name:      "percent encoding error in credentials",
			endpoint:  "https://%zz:" + secPassword + "@minio.internal:9000",
			bucketURL: "",
		},
		{
			name:      "single slash https with credentials in endpoint",
			endpoint:  "https:/admin:" + secPassword + "@minio.internal:9000",
			bucketURL: "",
		},
		{
			name:      "single slash http with credentials and port",
			endpoint:  "http:/admin:" + secPassword + "@minio:9000",
			bucketURL: "",
		},
		{
			name:      "single slash s3 with credentials in bucket url",
			endpoint:  "",
			bucketURL: "s3:/user:" + secPassword + "@host/x",
		},
		{
			name:      "secret in port position with absolute endpoint and credentials",
			endpoint:  "https://admin:" + secPassword + "@minio.internal:" + secPassword,
			bucketURL: "",
		},
		{
			name:      "secret in port position with absolute endpoint without credentials",
			endpoint:  "https://host:" + secPassword,
			bucketURL: "",
		},
		{
			name:      "secret in port position with schemeless endpoint",
			endpoint:  "minio.internal:" + secPassword,
			bucketURL: "",
		},
	}

	for _, tc := range malformedCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, storageLogModes()[0].debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = "test-bucket"
			cfg.Region = "us-east-1"
			cfg.AccessKeyID = "test-key"
			cfg.SecretAccessKey = secSecretKey
			cfg.Endpoint = tc.endpoint
			cfg.BucketURL = tc.bucketURL

			// ConnectWithConfig must reject malformed/whitespace inputs
			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.Error(t, err)

			// Assert that returned error message contains no secrets or credentials
			errStr := err.Error()
			assert.NotContains(t, errStr, secSecretKey)
			assert.NotContains(t, errStr, secPassword)
			assert.NotContains(t, errStr, secToken)

			// Assert direct config validation also fails safely without leaking secrets
			valErr := cfg.Validate()
			require.Error(t, valErr)
			assert.NotContains(t, valErr.Error(), secSecretKey)
			assert.NotContains(t, valErr.Error(), secPassword)
			assert.NotContains(t, valErr.Error(), secToken)

			// Assert log output contains no sensitive data
			rawOutput := buf.String()
			assertNoSensitiveDataInOutput(t, rawOutput)
		})
	}
}

func TestConnect_LegitimateEndpoints_AcceptedWithoutLeakingUserinfo(t *testing.T) {
	testCases := []struct {
		name       string
		endpoint   string
		expectHost string
	}{
		{
			name:       "ipv6 loopback with credentials and query",
			endpoint:   "http://admin:" + secPassword + "@[::1]:9000/data?token=" + secToken + "#" + secFragment,
			expectHost: "::1",
		},
		{
			name:       "schemeless ipv6 loopback with credentials",
			endpoint:   "admin:" + secPassword + "@[::1]:9000?token=" + secToken,
			expectHost: "::1",
		},
		{
			name:       "ipv4-mapped ipv6 with credentials",
			endpoint:   "http://admin:" + secPassword + "@[::ffff:127.0.0.1]:9000?token=" + secToken,
			expectHost: "::ffff:127.0.0.1",
		},
		{
			name:       "plain ipv6 loopback",
			endpoint:   "http://[::1]:9000",
			expectHost: "::1",
		},
		{
			name:       "at sign in reverse-proxy path prefix with credentials",
			endpoint:   "https://admin:" + secPassword + "@minio.local/tenant@team?token=" + secToken + "#" + secFragment,
			expectHost: "minio.local",
		},
		{
			name:       "percent-encoded at sign in path with credentials",
			endpoint:   "https://admin:" + secPassword + "@minio.local/a%40b?token=" + secToken,
			expectHost: "minio.local",
		},
		{
			name:       "at sign in path without credentials",
			endpoint:   "https://proxy.example.com/s3/a@b/",
			expectHost: "proxy.example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, storageLogModes()[0].debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = "test-bucket"
			cfg.Region = "us-east-1"
			cfg.AccessKeyID = "test-access-key-id"
			cfg.SecretAccessKey = secSecretKey
			cfg.SessionToken = secSessionToken
			cfg.Endpoint = tc.endpoint

			require.NoError(t, cfg.Validate())
			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.NoError(t, err)

			assertNoSensitiveDataInOutput(t, buf.String())

			detectLog := findStorageLog(parseStorageLogs(t, buf.Bytes()), "evaluating storage provider detection")
			require.NotNil(t, detectLog)
			assert.Equal(t, tc.expectHost, detectLog["inspected_host"])
			assert.Equal(t, "s3-compatible", detectLog["detected_provider"])
		})
	}
}
