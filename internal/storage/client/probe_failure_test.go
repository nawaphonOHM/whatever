package client_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	shortPingTimeout = 100 * time.Millisecond
	generousTimeout  = 30 * time.Second
	listTypeParam    = "list-type"
)

func blockUntilDone(req *http.Request) (*http.Response, error) {
	<-req.Context().Done()
	return nil, req.Context().Err()
}

func newPingClient(rec *recorder) *client.Client {
	return client.NewClient(client.ClientParams{
		Raw:            buildRawClientForTest(rec.httpClient()),
		DefaultBucket:  "my-bucket",
		PingTimeout:    shortPingTimeout,
		RequestTimeout: generousTimeout,
	})
}

func TestPing_UsesPingTimeoutNotRequestTimeout(t *testing.T) {
	t.Parallel()

	c := client.NewClient(client.ClientParams{
		Raw:            buildRawClientForTest(newMockHTTPClient(blockUntilDone)),
		DefaultBucket:  "my-bucket",
		PingTimeout:    shortPingTimeout,
		RequestTimeout: generousTimeout,
	})
	start := time.Now()
	err := c.Ping(t.Context())
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestPing_DeniedHeadBucketFallsBackToList(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: func(req recordedRequest) *http.Response {
		if req.Method == http.MethodHead {
			return newHTTPResponse(http.StatusForbidden, "", nil)
		}
		return newHTTPResponse(http.StatusOK, listXML(false, "", ""), nil)
	}}
	require.NoError(t, newPingClient(rec).Ping(t.Context()))
	reqs := rec.all()
	require.Len(t, reqs, 2)
	assert.Equal(t, "2", reqs[1].Query.Get(listTypeParam))
	assert.Equal(t, "1", reqs[1].Query.Get("max-keys"))
}

func TestPing_FailureMapping(t *testing.T) {
	t.Parallel()

	denied := &recorder{respond: statusResponder(http.StatusForbidden)}
	assert.ErrorIs(t, newPingClient(denied).Ping(t.Context()), contracts.ErrAccessDenied)

	failing := &recorder{respond: statusResponder(http.StatusInternalServerError)}
	err := newPingClient(failing).Ping(t.Context())
	require.Error(t, err)
	assert.NotErrorIs(t, err, contracts.ErrBucketNotFound)
	assert.NotErrorIs(t, err, contracts.ErrAccessDenied)
}

func probeConfig() *config.Config {
	cfg := endpointConfig(testHTTPSEndpoint)
	cfg.PingTimeout = shortPingTimeout
	return cfg
}

func connectWithProbe(t *testing.T, fn func(*http.Request) (*http.Response, error)) error {
	t.Helper()
	_, err := client.ConnectWithConfig(t.Context(), probeConfig(),
		client.WithHTTPClient(newMockHTTPClient(fn)), client.WithDriverOptions(disableRetries))
	return err
}

// Not parallel: it swaps the process-wide exit hook of the central logger worker.
func TestConnect_StartupProbeFailureExitsAbnormally(t *testing.T) {
	var exitCode atomic.Int32
	prev := client.SetExitFunc(func(code int) { exitCode.Store(int32(code)) })
	t.Cleanup(func() { _ = client.SetExitFunc(prev) })

	err := connectWithProbe(t, func(*http.Request) (*http.Response, error) {
		return newHTTPResponse(http.StatusNotFound, "", nil), nil
	})
	require.ErrorIs(t, err, contracts.ErrBucketNotFound)
	assert.NotZero(t, exitCode.Load())

	exitCode.Store(0)
	err = connectWithProbe(t, blockUntilDone)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.NotZero(t, exitCode.Load())
}
