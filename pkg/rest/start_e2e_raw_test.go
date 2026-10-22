package rest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRawE2EBlueprint() *BluePrint {
	regHook := &RRestAPIRegistration{
		Group: APIGroupRaw, Prefix: "/v1/custom",
		Apis: []*ExportableAPI{{
			Path: "/hook", Method: GET,
			Handler: func(Context) Response { return OK("hook-payload") },
		}},
	}
	regEvents := &RRestAPIRegistration{
		Group: APIGroupRaw, Prefix: "/webhooks",
		Apis: []*ExportableAPI{{
			Path: "/events", Method: GET,
			Handler: func(Context) Response { return OK("events-payload") },
		}},
	}
	regPing := &RRestAPIRegistration{
		Group: APIGroupRaw, Prefix: "/ping",
		Apis: []*ExportableAPI{{
			Path: "", Method: GET,
			Handler: func(Context) Response { return OK("pong") },
		}},
	}
	regUsers := &RRestAPIRegistration{
		Group: APIGroupRaw, Prefix: "/users",
		Apis: []*ExportableAPI{{
			Path: "/:id", Method: GET,
			Handler: func(c Context) Response { return OK(c.Param("id")) },
		}},
	}
	return NewBluePrint().WithAPIs(regHook, regEvents, regPing, regUsers)
}

func assertCustomAndWebhookRoutes(t *testing.T, baseURL string) {
	t.Helper()
	status, body := httpGetBody(t, baseURL+"/v1/custom/hook")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "hook-payload")

	status, body = httpGetBody(t, baseURL+"/webhooks/events")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "events-payload")
}

func assertPingAndParamRoutes(t *testing.T, baseURL string) {
	t.Helper()
	status, body := httpGetBody(t, baseURL+"/ping")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "pong")

	status, body = httpGetBody(t, baseURL+"/users/alice")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "alice")
}

func assertRawProbesAndMissingPaths(t *testing.T, baseURL string) {
	t.Helper()
	assertStatusOK(t, baseURL+"/health")
	assertStatusOK(t, baseURL+"/ready")
	assertStatusNotFound(t, baseURL+"/api/v1/custom/hook")
	assertStatusNotFound(t, baseURL+"/api/ping")
	assertStatusNotFound(t, baseURL+"/api/users/alice")
}

func TestStartREST_RawGroupE2E(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, testGroupRaw)

	errCh := startRESTAsync(newRawE2EBlueprint())
	baseURL := fmt.Sprintf("http://%s:%d", testHost, port)

	assertCustomAndWebhookRoutes(t, baseURL)
	assertPingAndParamRoutes(t, baseURL)
	assertRawProbesAndMissingPaths(t, baseURL)

	signalAndWait(t, errCh)
}
