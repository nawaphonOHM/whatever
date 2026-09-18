package server

import (
	"context"
	"net"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	srv := New(&Config{
		Host: testHost,
		Port: port,
		Mode: gin.TestMode,
	})

	// Act
	err := srv.Start(context.Background())

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "http server failed to start")
}
