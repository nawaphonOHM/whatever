package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

// occupyPort listens on an ephemeral port and returns it.
func occupyPort(t *testing.T) (net.Listener, int) {
	t.Helper()
	l, err := net.Listen("tcp", testHost+":0")
	require.NoError(t, err)
	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return l, tcpAddr.Port
}

// TestServer_StartPortInUse fails when the bind port is taken.
func TestServer_StartPortInUse(t *testing.T) {
	// Arrange
	l, port := occupyPort(t)
	defer func() { assert.NoError(t, l.Close()) }()
	srv, err := New(&Config{
		Host: testHost,
		Port: port,
		Mode: gin.TestMode,
	})
	require.NoError(t, err)

	// Act
	err = srv.Start(context.Background())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "http server failed to start")
}

// startInBackground runs Start and returns its error channel.
func startInBackground(srv *Server) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(context.Background())
	}()
	time.Sleep(startWait)
	return errCh
}

// TestServer_Shutdown stops a running server via Shutdown.
func TestServer_Shutdown(t *testing.T) {
	// Arrange
	srv, _ := newBoundServer(t)
	errCh := startInBackground(srv)

	// Act
	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownWait,
	)
	defer cancel()
	err := srv.Shutdown(ctx)

	// Assert
	require.NoError(t, err)
	assert.NoError(t, waitErr(t, errCh))
}
