package storage_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

var (
	errReadFailure  = errors.New("read failure")
	errCloseFailure = errors.New("close failure")
)

type trackingReadCloser struct {
	reader   io.Reader
	closeErr error
	closed   bool
}

func (t *trackingReadCloser) Read(p []byte) (int, error) {
	return t.reader.Read(p)
}

func (t *trackingReadCloser) Close() error {
	t.closed = true
	return t.closeErr
}

type errReader struct{}

func (*errReader) Read([]byte) (int, error) {
	return 0, errReadFailure
}

func TestReadAllBody_SuccessAndClosure(t *testing.T) {
	t.Parallel()

	content := []byte("stored file content")
	tracker := &trackingReadCloser{reader: bytes.NewReader(content)}
	out := &storage.GetObjectOutput{Body: tracker}

	data, err := storage.ReadAllBody(out)
	require.NoError(t, err)
	assert.Equal(t, content, data)
	assert.True(t, tracker.closed)
}

func TestReadAllBody_Errors(t *testing.T) {
	t.Parallel()

	t.Run("read error closes body and propagates error", func(t *testing.T) {
		t.Parallel()
		tracker := &trackingReadCloser{reader: &errReader{}}
		out := &storage.GetObjectOutput{Body: tracker}

		data, err := storage.ReadAllBody(out)
		assert.ErrorIs(t, err, errReadFailure)
		assert.Empty(t, data)
		assert.True(t, tracker.closed)
	})

	t.Run("close error propagates error", func(t *testing.T) {
		t.Parallel()
		tracker := &trackingReadCloser{
			reader:   bytes.NewReader([]byte("sample")),
			closeErr: errCloseFailure,
		}
		out := &storage.GetObjectOutput{Body: tracker}

		data, err := storage.ReadAllBody(out)
		assert.ErrorIs(t, err, errCloseFailure)
		assert.Equal(t, []byte("sample"), data)
		assert.True(t, tracker.closed)
	})
}

func TestReadAllBody_NilHandling(t *testing.T) {
	t.Parallel()

	dataNilOut, errNilOut := storage.ReadAllBody(nil)
	assert.NoError(t, errNilOut)
	assert.Nil(t, dataNilOut)

	dataNilBody, errNilBody := storage.ReadAllBody(&storage.GetObjectOutput{Body: nil})
	assert.NoError(t, errNilBody)
	assert.Nil(t, dataNilBody)
}
