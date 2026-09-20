package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/internal/rest/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyRegs returns an empty registration slice.
func emptyBlueprint() *contracts.BluePrint {
	return contracts.NewBluePrint()
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
	defer func() { assert.NoError(t, resp.Body.Close()) }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body testSuccessResponse[health.Status]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body
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
