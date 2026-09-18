package server

import (
	"fmt"
	"testing"
)

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

	// Shutdown
	signalAndWait(t, errCh)
}
