package server

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Shared test constants keep environment keys and timings consistent.
const (
	testHost        = "127.0.0.1"
	startWait       = 50 * time.Millisecond
	shutdownWait    = 3 * time.Second
	envGinMode      = "OHM9969_GIN_MODE"
	envServerPort   = "OHM9969_SERVER_PORT"
	envServerHost   = "OHM9969_SERVER_HOST"
	envAppVersion   = "OHM9969_APP_VERSION"
	customPort      = 9090
	customReadSec   = 5
	customIdleSec   = 30
	appVersionValue = "2.0.0"
)

// getFreePort reserves an ephemeral local TCP port for tests.
func getFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", testHost+":0")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, l.Close())
	}()
	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return tcpAddr.Port
}

// newBoundServer builds a test server on a free port with /ping.
func newBoundServer(t *testing.T) (*Server, int) {
	t.Helper()
	port := getFreePort(t)
	srv := New(&Config{
		Host: testHost,
		Port: port,
		Mode: gin.TestMode,
	})
	// /ping is the shared smoke endpoint for lifecycle tests.
	srv.Engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return srv, port
}

// waitErr waits for an error channel with timeout.
func waitErr(t *testing.T, errCh <-chan error) error {
	t.Helper()
	select {
	case err := <-errCh:
		return err
	case <-time.After(shutdownWait):
		t.Fatal("timed out waiting for server exit")
		return nil
	}
}

// setStartEnv configures host/port/mode for StartREST tests.
func setStartEnv(t *testing.T, port int) {
	t.Helper()
	// OHM9969_ keys match production config loading.
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envServerHost, testHost)
	t.Setenv(envServerPort, fmt.Sprintf("%d", port))
	t.Setenv(envAppVersion, appVersionValue)
}
