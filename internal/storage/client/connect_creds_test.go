package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func isolateAWSEnv(t *testing.T) {
	t.Helper()
	t.Setenv("AWS_CONFIG_FILE", "/nonexistent/path/config")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/nonexistent/path/credentials")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_ENDPOINT_URL", "")
	t.Setenv("AWS_ENDPOINT_URL_S3", "")
}

func setAWSEnvCredentials(t *testing.T) {
	t.Helper()
	isolateAWSEnv(t)
	t.Setenv("AWS_ACCESS_KEY_ID", testAccessKey)
	t.Setenv("AWS_SECRET_ACCESS_KEY", testSecretKey)
}

func TestConnect_DefaultCredentialChain_Signs(t *testing.T) {
	setAWSEnvCredentials(t)
	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	cfg.Region = testAWSRegionSE
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, testAWSRegionalHost, testVirtualPath)
	assert.Contains(t, got.Auth, testAccessKey)
}

func TestConnect_DefaultCredentialChain_Presign(t *testing.T) {
	setAWSEnvCredentials(t)
	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	httpClient := newMockHTTPClient(func(*http.Request) (*http.Response, error) {
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})
	c, err := client.ConnectWithConfig(
		context.Background(),
		cfg,
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
	)
	require.NoError(t, err)
	out, err := c.PresignPutObject(context.Background(), contracts.PresignPutObjectInput{
		Key:         "upload.png",
		ContentType: "image/png",
		ExpiresIn:   testPresignExpiry15m,
	})
	require.NoError(t, err)
	assert.Contains(t, out.URL, "upload.png")
	assert.Contains(t, out.URL, sigV4Algorithm)
}

func TestConnect_LoadDefaultConfig_InvalidProfile(t *testing.T) {
	isolateAWSEnv(t)
	t.Setenv("AWS_PROFILE", "doesnotexist_profile_xyz")
	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	httpClient := newMockHTTPClient(func(*http.Request) (*http.Response, error) {
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})
	_, err := client.ConnectWithConfig(
		context.Background(),
		cfg,
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load AWS default config")
}

func TestConnect_StaticKeysOverrideEnvCredentials(t *testing.T) {
	isolateAWSEnv(t)
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAENVCREDENTIALS")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "env-secret")
	cfg := testConfigWithKeys()
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Contains(t, got.Auth, testAccessKey)
	assert.NotContains(t, got.Auth, "AKIAENVCREDENTIALS")
}
