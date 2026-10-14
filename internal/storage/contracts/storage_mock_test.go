package contracts_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type mockStorage struct{}

func (*mockStorage) PutObject(context.Context, contracts.PutObjectInput) (*contracts.PutObjectOutput, error) {
	return &contracts.PutObjectOutput{ETag: "\"mock\""}, nil
}

func (*mockStorage) GetObject(context.Context, contracts.GetObjectInput) (*contracts.GetObjectOutput, error) {
	return &contracts.GetObjectOutput{Body: io.NopCloser(bytes.NewReader(nil))}, nil
}

func (*mockStorage) DeleteObject(context.Context, contracts.DeleteObjectInput) (*contracts.DeleteObjectOutput, error) {
	return &contracts.DeleteObjectOutput{}, nil
}

func (*mockStorage) HeadObject(context.Context, contracts.HeadObjectInput) (*contracts.HeadObjectOutput, error) {
	return &contracts.HeadObjectOutput{}, nil
}

func (*mockStorage) ListObjects(context.Context, contracts.ListObjectsInput) (*contracts.ListObjectsOutput, error) {
	return &contracts.ListObjectsOutput{}, nil
}

func (*mockStorage) Exists(ctx context.Context, bucket, key string) (bool, error) {
	_ = ctx
	return bucket != "" && key != "", nil
}

func (*mockStorage) PresignGetObject(
	context.Context,
	contracts.PresignGetObjectInput,
) (*contracts.PresignedURL, error) {
	return &contracts.PresignedURL{URL: "https://example.com/get"}, nil
}

func (*mockStorage) PresignPutObject(
	context.Context,
	contracts.PresignPutObjectInput,
) (*contracts.PresignedURL, error) {
	return &contracts.PresignedURL{URL: "https://example.com/put"}, nil
}

func (*mockStorage) Ping(context.Context) error {
	return nil
}

func (*mockStorage) Provider() contracts.Provider {
	return contracts.ProviderAWS
}

func (*mockStorage) DefaultBucket() string {
	return "test-bucket"
}

func (*mockStorage) RawClient() *s3.Client {
	return nil
}

func TestStorage_BasicMethods(t *testing.T) {
	t.Parallel()
	var s contracts.Storage = &mockStorage{}
	ctx := context.Background()

	assert.Equal(t, contracts.ProviderAWS, s.Provider())
	assert.Equal(t, "test-bucket", s.DefaultBucket())
	assert.Nil(t, s.RawClient())
	assert.NoError(t, s.Ping(ctx))

	exists, err := s.Exists(ctx, "b", "k")
	assert.NoError(t, err)
	assert.True(t, exists)
}
