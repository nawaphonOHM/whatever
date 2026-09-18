package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertPingOK checks the /ping endpoint body.
func assertPingOK(t *testing.T, port int) {
	t.Helper()
	url := fmt.Sprintf("http://%s:%d/ping", testHost, port)
	resp, err := http.Get(url)
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))
}

// startWithCancel runs Start under a cancelable context.
func startWithCancel(srv *Server) (context.CancelFunc, <-chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start(ctx) }()
	time.Sleep(startWait)
	return cancel, errCh
}

// TestServer_StartAndGracefulShutdown serves then cancels cleanly.
func TestServer_StartAndGracefulShutdown(t *testing.T) {
	// Arrange
	srv, port := newBoundServer(t)
	cancel, errCh := startWithCancel(srv)

	// Act / Assert request path
	assertPingOK(t, port)

	// Shutdown via cancel
	cancel()
	assert.NoError(t, waitErr(t, errCh))
}
