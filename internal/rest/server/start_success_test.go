package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyBlueprint returns an empty registration blueprint.
func emptyBlueprint() *contracts.BluePrint {
	return contracts.NewBluePrint()
}

// pingBlueprint returns a versioned /ping blueprint.
func pingBlueprint() *contracts.BluePrint {
	return contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
		Version: 1,
		Prefix:  "/api",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/ping",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("pong")
			},
		}},
	})
}

// startRESTAsync launches StartREST and waits for bind.
func startRESTAsync(
	bluePrint *contracts.BluePrint,
) <-chan error {
	errCh := make(chan error, 1)
	go func() { errCh <- StartREST(bluePrint) }()
	time.Sleep(startWait)
	return errCh
}

// signalAndWait sends SIGTERM and waits for StartREST exit.
func signalAndWait(t *testing.T, errCh <-chan error) {
	t.Helper()
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
	assert.NoError(t, waitErr(t, errCh))
}

func closeRespBody(t *testing.T, resp *http.Response) {
	t.Helper()
	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		t.Logf("failed to discard body: %v", copyErr)
	}
	assert.NoError(t, resp.Body.Close())
	http.DefaultClient.CloseIdleConnections()
}

// assertStatusOK performs GET and asserts status OK.
func assertStatusOK(t *testing.T, url string) {
	t.Helper()
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer closeRespBody(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// fetchHealth decodes the reserved health payload.
func fetchHealth(
	t *testing.T,
	port int,
) testSuccessResponse[health.Status] {
	t.Helper()
	url := fmt.Sprintf(
		"http://%s:%d%s", testHost, port, ReservedHealthPath,
	)
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer closeRespBody(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body testSuccessResponse[health.Status]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body
}

// TestStartREST_Success boots, serves /api/v1/ping, then exits.
func TestStartREST_Success(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(pingBlueprint())

	// Act
	url := fmt.Sprintf(
		"http://%s:%d/api/v1/ping", testHost, port,
	)
	assertStatusOK(t, url)

	// Shutdown via signal
	signalAndWait(t, errCh)
}

// TestStartREST_HealthEndpoint serves reserved /health with version.
func TestStartREST_HealthEndpoint(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(emptyBlueprint())

	// Act
	body := fetchHealth(t, port)

	// Assert
	assert.Equal(t, appVersionValue, body.Data.Version)
	signalAndWait(t, errCh)
}
