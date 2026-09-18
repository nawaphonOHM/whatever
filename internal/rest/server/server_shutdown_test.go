package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
