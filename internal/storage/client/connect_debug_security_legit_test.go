package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func legitimateEndpointCases() []struct {
	name       string
	endpoint   string
	expectHost string
} {
	return []struct {
		name       string
		endpoint   string
		expectHost string
	}{
		{
			name:       "ipv6 loopback with credentials and query",
			endpoint:   "http://admin:" + secPassword + "@[::1]:9000/data?token=" + secToken + secHashFrag,
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
			endpoint:   secAdminPrefix + secPassword + "@minio.local/tenant@team?token=" + secToken + secHashFrag,
			expectHost: "minio.local",
		},
		{
			name:       "percent-encoded at sign in path with credentials",
			endpoint:   secAdminPrefix + secPassword + "@minio.local/a%40b?token=" + secToken,
			expectHost: "minio.local",
		},
		{
			name:       "at sign in path without credentials",
			endpoint:   "https://proxy.example.com/s3/a@b/",
			expectHost: "proxy.example.com",
		},
	}
}

func TestConnect_LegitimateEndpoints_AcceptedWithoutLeakingUserinfo(t *testing.T) {
	for _, tc := range legitimateEndpointCases() {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, storageLogModes()[0].debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = secTestBucket
			cfg.Region = secTestRegion
			cfg.AccessKeyID = secTestAccessKey
			cfg.SecretAccessKey = secSecretKey
			cfg.SessionToken = secSessionToken
			cfg.Endpoint = tc.endpoint

			require.NoError(t, cfg.Validate())
			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.NoError(t, err)

			assertNoSensitiveDataInOutput(t, buf.String())

			detectLog := findStorageLog(parseStorageLogs(t, buf.Bytes()), secLogEvalDetect)
			require.NotNil(t, detectLog)
			assert.Equal(t, tc.expectHost, detectLog[secLogInspectHost])
			assert.Equal(t, secProviderS3, detectLog[secLogDetectProv])
		})
	}
}
