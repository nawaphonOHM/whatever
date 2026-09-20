package rest

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helpers for E2E tests.

const (
	testHost      = "127.0.0.1"
	startWait     = 50 * time.Millisecond
	shutdownWait  = 3 * time.Second
	envGinMode    = "OHM9996_GIN_MODE"
	envServerPort = "OHM9996_SERVER_PORT"
	envServerHost = "OHM9996_SERVER_HOST"
	envAppVersion = "OHM9996_APP_VERSION"
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

// setStartEnv configures host/port/mode for StartREST tests.
func setStartEnv(t *testing.T, port int) {
	t.Helper()
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envServerHost, testHost)
	t.Setenv(envServerPort, fmt.Sprintf("%d", port))
	t.Setenv(envAppVersion, "1.0.0")
}

// startRESTAsync launches StartREST and waits for bind.
func startRESTAsync(
	bluePrint *BluePrint,
) <-chan error {
	errCh := make(chan error, 1)
	go func() { errCh <- StartREST(bluePrint) }()
	time.Sleep(startWait)
	return errCh
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

// signalAndWait sends SIGTERM and waits for StartREST exit.
func signalAndWait(t *testing.T, errCh <-chan error) {
	t.Helper()
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
	assert.NoError(t, waitErr(t, errCh))
}

// assertStatusOK performs GET and asserts status OK.
func assertStatusOK(t *testing.T, url string) {
	t.Helper()
	resp, err := http.Get(url)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, resp.Body.Close())
}
