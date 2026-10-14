package contracts_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestStorage_MutatingOperations(t *testing.T) {
	t.Parallel()
	var s contracts.Storage = &mockStorage{}
	ctx := context.Background()

	putOut, err := s.PutObject(ctx, contracts.PutObjectInput{})
	assert.NoError(t, err)
	assert.Equal(t, "\"mock\"", putOut.ETag)

	delOut, err := s.DeleteObject(ctx, contracts.DeleteObjectInput{})
	assert.NoError(t, err)
	assert.NotNil(t, delOut)
}

func TestStorage_ReadOperations(t *testing.T) {
	t.Parallel()
	var s contracts.Storage = &mockStorage{}
	ctx := context.Background()

	getOut, err := s.GetObject(ctx, contracts.GetObjectInput{})
	assert.NoError(t, err)
	require.NotNil(t, getOut.Body)
	assert.NoError(t, getOut.Body.Close())

	headOut, err := s.HeadObject(ctx, contracts.HeadObjectInput{})
	assert.NoError(t, err)
	assert.NotNil(t, headOut)
}

func TestStorage_List(t *testing.T) {
	t.Parallel()
	var s contracts.Storage = &mockStorage{}
	ctx := context.Background()

	listOut, err := s.ListObjects(ctx, contracts.ListObjectsInput{})
	assert.NoError(t, err)
	assert.NotNil(t, listOut)
}

func TestStorage_Presign(t *testing.T) {
	t.Parallel()
	var s contracts.Storage = &mockStorage{}
	ctx := context.Background()

	presignGet, err := s.PresignGetObject(ctx, contracts.PresignGetObjectInput{})
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/get", presignGet.URL)

	presignPut, err := s.PresignPutObject(ctx, contracts.PresignPutObjectInput{})
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/put", presignPut.URL)
}
