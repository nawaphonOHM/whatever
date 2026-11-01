package client_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func TestConnect_Checksum_AWS_DefaultChain(t *testing.T) {
	setAWSEnvCredentials(t)

	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	cfg.Region = "us-east-1"
	c := connectForChecksum(t, cfg)
	assertWhenSupportedChecksum(t, c.RawClient().Options())
}

func TestConnect_Checksum_AWS_ExplicitEnvRespected(t *testing.T) {
	setAWSEnvCredentials(t)
	t.Setenv("AWS_REQUEST_CHECKSUM_CALCULATION", "WHEN_REQUIRED")
	t.Setenv("AWS_RESPONSE_CHECKSUM_VALIDATION", "WHEN_REQUIRED")

	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	cfg.Region = "us-east-1"
	c := connectForChecksum(t, cfg)
	assertWhenRequiredChecksum(t, c.RawClient().Options())
}

func TestConnect_Checksum_S3Compatible_EnvOverridesIgnored(t *testing.T) {
	setOpposingChecksumEnv(t)

	cfgStatic := testConfigWithKeys()
	cfgStatic.Endpoint = "https://minio.internal.net"
	assertConnectChecksumWhenRequired(t, cfgStatic)

	cfgDefault := config.DefaultConfig()
	cfgDefault.Bucket = testBucket
	cfgDefault.Endpoint = "https://minio.internal.net"
	assertConnectChecksumWhenRequired(t, cfgDefault)
}
