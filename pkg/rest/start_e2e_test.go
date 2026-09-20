package rest

import (
	"fmt"
	"testing"
)

// TestStartREST_Success boots, serves a registered route, then exits.
func TestStartREST_Success(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)

	api := &ExportableAPI{
		Path:   "/ping",
		Method: GET,
		Handler: func(Context) Response {
			return OK("pong")
		},
	}
	reg := &RRestAPIRegistration{
		Version: 1,
		Prefix:  "/api",
		Apis:    []*ExportableAPI{api},
	}
	bp := NewBluePrint().WithAPIs(reg)
	errCh := startRESTAsync(bp)

	// Act
	url := fmt.Sprintf(
		"http://%s:%d/v1/api/ping", testHost, port,
	)
	assertStatusOK(t, url)

	// Shutdown via signal
	signalAndWait(t, errCh)
}

// TestStartREST_CustomCORS verifies custom CORS headers are applied.
func TestStartREST_CustomCORS(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	bp := newCustomCORSBlueprint()
	errCh := startRESTAsync(bp)

	// Act: Send OPTIONS preflight request
	url := fmt.Sprintf("http://%s:%d/health", testHost, port)
	resp := corsPreflightRequest(t, url)
	defer closeBody(t, resp.Body)

	// Assert CORS headers match custom configuration
	assertCORSHeaders(t, resp, "https://example.com", "POST, PATCH")

	// Shutdown via signal
	signalAndWait(t, errCh)
}

// TestStartREST_DefaultCORS verifies default CORS headers when unconfigured.
func TestStartREST_DefaultCORS(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	bp := NewBluePrint()
	errCh := startRESTAsync(bp)

	// Act: Send OPTIONS preflight request
	url := fmt.Sprintf("http://%s:%d/health", testHost, port)
	resp := corsPreflightRequest(t, url)
	defer closeBody(t, resp.Body)

	// Assert default permissive CORS headers
	assertCORSHeaders(t, resp, "*", "")

	// Shutdown via signal
	signalAndWait(t, errCh)
}
