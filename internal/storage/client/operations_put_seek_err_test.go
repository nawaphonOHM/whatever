package client_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const seekStartOffset = 6

var errMockSeek = errors.New("mock seek error")

type failingSeeker struct {
	failWhence int
	pos        int64
}

func (*failingSeeker) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (s *failingSeeker) Seek(offset int64, whence int) (int64, error) {
	if whence == s.failWhence {
		return 0, errMockSeek
	}
	s.pos = offset
	return s.pos, nil
}

func TestPutObject_SeekableErrorBranches(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: statusResponder(http.StatusOK)}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)

	t.Run("FailingSeekEnd returns serialization error", func(t *testing.T) {
		body := &failingSeeker{failWhence: io.SeekEnd}
		_, err := c.PutObject(t.Context(), contracts.PutObjectInput{Key: testObjectKey, Body: body})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mock seek error")
	})

	t.Run("FailingSeekStart returns serialization error", func(t *testing.T) {
		body := &failingSeeker{failWhence: io.SeekStart}
		_, err := c.PutObject(t.Context(), contracts.PutObjectInput{Key: testObjectKey, Body: body})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mock seek error")
	})
}

func assertSingleBody(t *testing.T, rec *recorder, expected string) {
	t.Helper()
	reqs := rec.all()
	require.Len(t, reqs, 1)
	assert.Equal(t, expected, string(reqs[0].Body))
}

func TestPutObject_SeekPositionRestored(t *testing.T) {
	t.Parallel()

	reader := strings.NewReader("hello world from storage")
	_, err := reader.Seek(seekStartOffset, io.SeekStart)
	require.NoError(t, err)

	rec := &recorder{respond: statusResponder(http.StatusOK)}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	_, err = c.PutObject(t.Context(), contracts.PutObjectInput{Key: testObjectKey, Body: reader})
	require.NoError(t, err)
	assertSingleBody(t, rec, "world from storage")
}
