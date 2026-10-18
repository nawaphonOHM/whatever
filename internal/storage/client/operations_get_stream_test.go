package client_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	shortRequestTimeout = 100 * time.Millisecond
	slowBodyDelay       = 300 * time.Millisecond
	maxStreamWait       = 5 * time.Second
)

// ctxBody emulates a network body: each Read waits delay unless the request context ends first.
type ctxBody struct {
	ctx    context.Context
	chunks []string
	delay  time.Duration
}

func (b *ctxBody) wait() error {
	select {
	case <-time.After(b.delay):
		return nil
	case <-b.ctx.Done():
		return b.ctx.Err()
	}
}

func (b *ctxBody) Read(p []byte) (int, error) {
	if len(b.chunks) == 0 {
		return 0, io.EOF
	}
	if err := b.wait(); err != nil {
		return 0, err
	}
	n := copy(p, b.chunks[0])
	b.chunks = b.chunks[1:]
	return n, nil
}

func (*ctxBody) Close() error { return nil }

func newGetClient(fn func(*http.Request) (*http.Response, error)) *client.Client {
	return client.NewClient(client.ClientParams{
		Raw:            buildRawClientForTest(newMockHTTPClient(fn)),
		DefaultBucket:  "my-bucket",
		RequestTimeout: shortRequestTimeout,
	})
}

func slowBodyResponse(req *http.Request) (*http.Response, error) {
	resp := newHTTPResponse(http.StatusOK, "", map[string]string{"Content-Length": "11"})
	resp.Body = &ctxBody{ctx: req.Context(), chunks: []string{"hello ", "world"}, delay: slowBodyDelay}
	return resp, nil
}

func TestGetObject_BodyStreamsPastRequestTimeout(t *testing.T) {
	t.Parallel()

	c := newGetClient(slowBodyResponse)
	out, err := c.GetObject(t.Context(), contracts.GetObjectInput{Key: testObjectKey})
	require.NoError(t, err)
	data, err := io.ReadAll(out.Body)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(data))
	require.NoError(t, out.Body.Close())
}

func TestGetObject_HeaderTimeout(t *testing.T) {
	t.Parallel()

	c := newGetClient(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, req.Context().Err()
	})
	start := time.Now()
	out, err := c.GetObject(t.Context(), contracts.GetObjectInput{Key: testObjectKey})
	assert.Nil(t, out)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, time.Since(start), maxStreamWait)
}

func TestGetObject_CallerCancelAbortsBody(t *testing.T) {
	t.Parallel()

	c := newGetClient(slowBodyResponse)
	ctx, cancel := context.WithCancel(t.Context())
	out, err := c.GetObject(ctx, contracts.GetObjectInput{Key: testObjectKey})
	require.NoError(t, err)
	time.AfterFunc(shortRequestTimeout, cancel)
	_, err = io.ReadAll(out.Body)
	assert.ErrorIs(t, err, context.Canceled)
	require.NoError(t, out.Body.Close())
}

func TestGetObject_CloseReleasesRequestContext(t *testing.T) {
	t.Parallel()

	var reqCtx context.Context
	c := newGetClient(func(req *http.Request) (*http.Response, error) {
		reqCtx = req.Context()
		return newHTTPResponse(http.StatusOK, "data", nil), nil
	})
	out, err := c.GetObject(t.Context(), contracts.GetObjectInput{Key: testObjectKey})
	require.NoError(t, err)
	require.NoError(t, reqCtx.Err())
	require.NoError(t, out.Body.Close())
	assert.Error(t, reqCtx.Err())
}
