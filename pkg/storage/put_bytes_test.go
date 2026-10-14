package storage_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestNewPutObjectBytes(t *testing.T) {
	t.Parallel()

	payload := []byte("payload data")
	input := storage.NewPutObjectBytes("my-bucket", "my/key.txt", payload, "text/plain")

	assert.Equal(t, "my-bucket", input.Bucket)
	assert.Equal(t, "my/key.txt", input.Key)
	assert.Equal(t, "text/plain", input.ContentType)
	require.NotNil(t, input.Body)

	readData, err := io.ReadAll(input.Body)
	require.NoError(t, err)
	assert.Equal(t, payload, readData)
}
