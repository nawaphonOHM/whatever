package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

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
			cfg.Bucket = secTestBucket
			cfg.Region = secTestRegion
			cfg.AccessKeyID = secTestAccessKey
			cfg.SecretAccessKey = secSecretKey
			cfg.SessionToken = secSessionToken
			cfg.Endpoint = "admin:" + secPassword + "@" + secHostMinio +
				":9000/data?redir=https://other.internal/x&token=" + secToken + secHashFrag
			cfg.BucketURL = secUserPrefix + secPassword + "@custom-host:9000/data?r=a://b&token=" +
				secToken + "#frag://b"

			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.NoError(t, err)

			rawOutput := buf.String()
			assertNoSensitiveDataInOutput(t, rawOutput)

			entries := parseStorageLogs(t, buf.Bytes())
			detectLog := findStorageLog(entries, secLogEvalDetect)
			require.NotNil(t, detectLog)
			assert.Equal(t, "minio.internal:9000/data", detectLog["raw_endpoint"])
			assert.Equal(t, "custom-host:9000/data", detectLog["raw_bucket_url"])
			assert.Equal(t, secHostMinio, detectLog[secLogInspectHost])
			assert.Equal(t, secProviderS3, detectLog[secLogDetectProv])
		})
	}
}
