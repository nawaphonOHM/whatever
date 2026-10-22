package rest

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envAPIGroup       = "OHM9996_SERVER_API_GROUP"
	testGroupRaw      = "raw"
	testGroupStd      = "standard"
	testCustomPath    = "/custom"
	errMismatchSubstr = "does not match OHM9996_SERVER_API_GROUP"
)

// assertResponseStatus executes a GET request and asserts the HTTP status code.
func assertResponseStatus(t *testing.T, url string, expectedStatus int) {
	t.Helper()
	resp, err := http.Get(url)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		t.Logf("failed to discard body: %v", copyErr)
	}
	assert.NoError(t, resp.Body.Close())
	http.DefaultClient.CloseIdleConnections()
}

// assertStatusNotFound asserts GET request returns 404 Not Found.
func assertStatusNotFound(t *testing.T, url string) {
	t.Helper()
	assertResponseStatus(t, url, http.StatusNotFound)
}

// httpGetBody performs a GET request and returns status and body string.
func httpGetBody(t *testing.T, url string) (int, string) {
	t.Helper()
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
		http.DefaultClient.CloseIdleConnections()
	}()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, string(body)
}

// httpPostBody performs a POST request and returns status and body string.
func httpPostBody(t *testing.T, url, payload string) (int, string) {
	t.Helper()
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(payload))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
		http.DefaultClient.CloseIdleConnections()
	}()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, string(body)
}

// guardedSIGTERM registers a local SIGTERM listener before signaling the process,
// so the test binary survives even when no StartREST listener is registered.
func guardedSIGTERM(t *testing.T) func() {
	t.Helper()
	guard := make(chan os.Signal, 1)
	signal.Notify(guard, syscall.SIGTERM)
	assert.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
	return func() { signal.Stop(guard) }
}

// cleanupTimeoutServer stops a regression server and then fails the test.
func cleanupTimeoutServer(t *testing.T, errCh <-chan error) {
	t.Helper()
	stopGuard := guardedSIGTERM(t)
	defer stopGuard()
	select {
	case err := <-errCh:
		t.Logf("regression server exited: %v", err)
	case <-time.After(shutdownWait):
		t.Log("timed out waiting for regression server cleanup")
	}
	t.Fatal("StartREST did not fail fast within timeout; server may have started")
}

// startRESTWithTimeout runs StartREST bounded by a short failure timeout.
func startRESTWithTimeout(t *testing.T, bp *BluePrint) error {
	t.Helper()
	setStartEnv(t, getFreePort(t))
	errCh := make(chan error, 1)
	go func() { errCh <- StartREST(bp) }()
	select {
	case err := <-errCh:
		return err
	case <-time.After(1 * time.Second):
		cleanupTimeoutServer(t, errCh)
		return nil
	}
}
