package server

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testHost     = "127.0.0.1"
	startWait    = 50 * time.Millisecond
	shutdownWait = 3 * time.Second
	envGinMode   = "OHM9969_GIN_MODE"
	envHost      = "OHM9969_SERVER_HOST"
	envPort      = "OHM9969_SERVER_PORT"
)

// getFreePort reserves an ephemeral local TCP port.
func getFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", testHost+":0")
	require.NoError(t, err)
	defer func() { assert.NoError(t, l.Close()) }()
	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return tcpAddr.Port
}

// pingRegs returns a versioned /ping registration set.
func pingRegs() []*RestAPIRegistration {
	return []*RestAPIRegistration{{
		Version: 1,
		Prefix:  "/api",
		Apis: []*ExportableAPI{{
			Path:   "/ping",
			Method: GET,
			Handler: func(*Context) response.Response {
				return response.OK("pong")
			},
		}},
	}}
}

// setStartEnv configures host/port/mode for StartREST tests.
func setStartEnv(t *testing.T, port int) {
	t.Helper()
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envHost, testHost)
	t.Setenv(envPort, fmt.Sprintf("%d", port))
}

// startRESTAsync launches StartREST and waits for bind.
func startRESTAsync(regs []*RestAPIRegistration) <-chan error {
	errCh := make(chan error, 1)
	go func() { errCh <- StartREST(regs) }()
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
