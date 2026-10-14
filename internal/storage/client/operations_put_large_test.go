package client_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	largeSeekerSize     = 6 * 1024 * 1024 * 1024
	oversizedSeekerSize = 100 * 1024 * 1024 * 1024
	fakeSeekerBytes     = 100
)

type mockLargeSeeker struct {
	reader io.Reader
	size   int64
	pos    int64
}

func (s *mockLargeSeeker) Read(p []byte) (int, error) {
	n, err := s.reader.Read(p)
	s.pos += int64(n)
	return n, err
}

func resolveSeekOffset(s *mockLargeSeeker, offset int64, whence int) int64 {
	if whence == io.SeekStart {
		return offset
	}
	if whence == io.SeekEnd {
		return s.size + offset
	}
	return s.pos + offset
}

func (s *mockLargeSeeker) Seek(offset int64, whence int) (int64, error) {
	s.pos = resolveSeekOffset(s, offset, whence)
	return s.pos, nil
}

func assertLargeSeekableResult(t *testing.T, rec *recorder, body *mockLargeSeeker, err error) {
	t.Helper()
	require.NoError(t, err)
	reqs := rec.all()
	require.Len(t, reqs, multipartRequestCount)
	assert.True(t, reqs[0].Query.Has("uploads"))
	assert.Positive(t, body.pos)
}

func TestPutObject_LargeSeekableRoutesToMultipart(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: multipartResponder("0")}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	largeBody := &mockLargeSeeker{
		reader: &zeroReader{left: partSize + 1},
		size:   largeSeekerSize,
	}
	out, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: largeBody,
	})
	assert.Equal(t, `"final-etag"`, out.ETag)
	assertLargeSeekableResult(t, rec, largeBody, err)
}

func TestPutObject_OversizedSeekablePreflightGuard(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: multipartResponder("0")}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	oversizedBody := &mockLargeSeeker{
		reader: &zeroReader{left: fakeSeekerBytes},
		size:   oversizedSeekerSize,
	}
	_, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: oversizedBody,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds the maximum number of multipart parts")
	assert.Empty(t, rec.all())
}
