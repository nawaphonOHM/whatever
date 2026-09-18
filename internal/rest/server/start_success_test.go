package server

import (
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// startRESTAsync launches StartREST and waits for bind.
func startRESTAsync(
	regs []*RestAPIRegistration,
) <-chan error {
	errCh := make(chan error, 1)
	go func() { errCh <- StartREST(regs) }()
	time.Sleep(startWait)
	return errCh
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

// TestStartREST_Success boots, serves /v1/api/ping, then exits.
func TestStartREST_Success(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(pingRegs())

	// Act
	url := fmt.Sprintf(
		"http://%s:%d/v1/api/ping", testHost, port,
	)
	assertStatusOK(t, url)

	// Shutdown via signal
	signalAndWait(t, errCh)
}
