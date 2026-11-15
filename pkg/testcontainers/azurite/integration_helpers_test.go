//go:build testcontainers

package azurite_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	containerTimeout = 2 * time.Minute
)

type s3ResourceNames struct {
	bucket string
	key    string
}

func newAzuriteStorageClient(
	t *testing.T,
	ctx context.Context,
	endpoint string,
	names s3ResourceNames,
) (*storage.Client, func()) {
	t.Helper()
	cfg := newAzuriteStorageConfig(endpoint, names)
	httpClient := &http.Client{Timeout: s3GatewayRequestLimit}
	bootstrap := connectAzuriteStorageWithoutProbe(t, ctx, cfg, httpClient)
	_, err := bootstrap.RawClient().CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(names.bucket)})
	require.NoError(t, err)
	markBucketDeleted := registerS3BucketCleanup(t, bootstrap.RawClient(), names.bucket)
	return connectAzuriteStorage(t, ctx, cfg, httpClient), markBucketDeleted
}

func newAzuriteStorageConfig(endpoint string, names s3ResourceNames) *storage.Config {
	cfg := storage.DefaultConfig()
	cfg.Bucket = names.bucket
	cfg.Region = s3ProxyRegion
	cfg.AccessKeyID = s3ProxyIdentity
	cfg.SecretAccessKey = s3ProxyCredential
	cfg.Endpoint = endpoint
	cfg.UsePathStyle = true
	cfg.Insecure = true
	return cfg
}

func connectAzuriteStorage(
	t *testing.T,
	ctx context.Context,
	cfg *storage.Config,
	httpClient *http.Client,
) *storage.Client {
	t.Helper()
	client, err := storage.ConnectWithConfig(ctx, cfg, storage.WithHTTPClient(httpClient))
	require.NoError(t, err)
	return client
}

func connectAzuriteStorageWithoutProbe(
	t *testing.T,
	ctx context.Context,
	cfg *storage.Config,
	httpClient *http.Client,
) *storage.Client {
	t.Helper()
	client, err := storage.ConnectWithConfig(
		ctx,
		cfg,
		storage.WithSkipProbe(true),
		storage.WithHTTPClient(httpClient),
	)
	require.NoError(t, err)
	return client
}

func newS3ResourceNames() s3ResourceNames {
	return s3ResourceNames{
		bucket: fmt.Sprintf("azurite-s3-lifecycle-%d", time.Now().UnixNano()),
		key:    "lifecycle-test.txt",
	}
}
