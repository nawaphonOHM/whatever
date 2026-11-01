package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func assertWhenRequiredChecksum(t *testing.T, opts s3.Options) {
	t.Helper()
	assert.Equal(t, aws.RequestChecksumCalculationWhenRequired, opts.RequestChecksumCalculation)
	assert.Equal(t, aws.ResponseChecksumValidationWhenRequired, opts.ResponseChecksumValidation)
	assert.True(t, opts.DisableLogOutputChecksumValidationSkipped)
}

func assertWhenSupportedChecksum(t *testing.T, opts s3.Options) {
	t.Helper()
	assert.Equal(t, aws.RequestChecksumCalculationWhenSupported, opts.RequestChecksumCalculation)
	assert.Equal(t, aws.ResponseChecksumValidationWhenSupported, opts.ResponseChecksumValidation)
	assert.True(t, opts.DisableLogOutputChecksumValidationSkipped)
}

func connectForChecksum(t *testing.T, cfg *config.Config) *client.Client {
	t.Helper()
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
	return c
}

func assertConnectChecksumWhenRequired(t *testing.T, cfg *config.Config) {
	t.Helper()
	c := connectForChecksum(t, cfg)
	assertWhenRequiredChecksum(t, c.RawClient().Options())
}

func setOpposingChecksumEnv(t *testing.T) {
	t.Helper()
	isolateAWSEnv(t)
	t.Setenv("AWS_REQUEST_CHECKSUM_CALCULATION", "WHEN_SUPPORTED")
	t.Setenv("AWS_RESPONSE_CHECKSUM_VALIDATION", "WHEN_SUPPORTED")
	t.Setenv("AWS_ACCESS_KEY_ID", testAccessKey)
	t.Setenv("AWS_SECRET_ACCESS_KEY", testSecretKey)
}
