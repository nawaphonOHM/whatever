package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/rest/health"
	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyRegs returns an empty registration slice.
func emptyRegs() []*rest.RRestAPIRegistration {
	return make([]*rest.RRestAPIRegistration, 0)
}

// fetchHealth decodes the reserved health payload.
func fetchHealth(
	t *testing.T,
	port int,
) rest.SuccessResponse[health.Status] {
	t.Helper()
	url := fmt.Sprintf(
		"http://%s:%d%s", testHost, port, ReservedHealthPath,
	)
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer func() { assert.NoError(t, resp.Body.Close()) }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var body rest.SuccessResponse[health.Status]
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body
}

// TestStartREST_HealthEndpoint serves reserved /health with version.
func TestStartREST_HealthEndpoint(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(emptyRegs())

	// Act
	body := fetchHealth(t, port)

	// Assert
	assert.Equal(t, appVersionValue, body.Data.Version)
	signalAndWait(t, errCh)
}
