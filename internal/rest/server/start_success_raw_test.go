package server

import (
	"fmt"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
)

// rawPingBlueprint returns an unversioned raw /custom/ping blueprint.
func rawPingBlueprint() *contracts.BluePrint {
	return contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
		Group:   contracts.APIGroupRaw,
		Version: 1,
		Prefix:  "/custom",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/ping",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("raw_pong")
			},
		}},
	})
}

// TestStartREST_RawSuccess boots in raw mode, serves /custom/ping, and exits.
func TestStartREST_RawSuccess(t *testing.T) {
	// Arrange
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(rawPingBlueprint())

	// Act & Assert business route without /api prefix
	url := fmt.Sprintf("http://%s:%d/custom/ping", testHost, port)
	assertStatusOK(t, url)

	// Act & Assert reserved health probe works in raw mode
	body := fetchHealth(t, port)
	assert.Equal(t, appVersionValue, body.Data.Version)

	// Shutdown via signal
	signalAndWait(t, errCh)
}
