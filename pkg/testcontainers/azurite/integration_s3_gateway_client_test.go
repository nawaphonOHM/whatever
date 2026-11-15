//go:build testcontainers

package azurite_test

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/pkg/testcontainers/azurite"
)

func newS3GatewayClient(endpoint string) *s3.Client {
	awsConfig := aws.Config{
		Region:      s3ProxyRegion,
		Credentials: credentials.NewStaticCredentialsProvider(s3ProxyIdentity, s3ProxyCredential, ""),
		HTTPClient:  &http.Client{Timeout: s3GatewayRequestLimit},
	}
	return s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
		options.UsePathStyle = true
	})
}

func azuriteStopper(emulator *azurite.Container, markRemoved func()) func(context.Context) error {
	return func(ctx context.Context) error {
		if err := emulator.RawContainer().Terminate(ctx); err != nil {
			return err
		}
		markRemoved()
		return nil
	}
}
