package rest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newStandardE2EBlueprint() *BluePrint {
	regV2 := &RRestAPIRegistration{
		Group: APIGroupStandard, Version: 2, Prefix: "/orders",
		Apis: []*ExportableAPI{{
			Path: "/summary", Method: GET,
			Handler: func(Context) Response { return OK("orders-v2-summary") },
		}},
	}
	regV1 := &RRestAPIRegistration{
		Group: APIGroupStandard, Version: 1, Prefix: "/items",
		Apis: []*ExportableAPI{{
			Path: "/:id", Method: GET,
			Handler: func(c Context) Response { return OK("item-" + c.Param("id")) },
		}},
	}
	regV0 := &RRestAPIRegistration{
		Group: APIGroupStandard, Version: 0, Prefix: "/ping",
		Apis: []*ExportableAPI{{
			Path: "", Method: GET,
			Handler: func(Context) Response { return OK("standard-pong") },
		}},
	}
	return NewBluePrint().WithAPIs(regV2, regV1, regV0)
}

func assertStandardRoutes(t *testing.T, baseURL string) {
	t.Helper()
	status, body := httpGetBody(t, baseURL+"/api/v2/orders/summary")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "orders-v2-summary")

	status, body = httpGetBody(t, baseURL+"/api/v1/items/99")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "item-99")

	status, body = httpGetBody(t, baseURL+"/api/ping")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "standard-pong")
}

func assertStandardProbesAndMissingPaths(t *testing.T, baseURL string) {
	t.Helper()
	assertStatusOK(t, baseURL+"/health")
	assertStatusOK(t, baseURL+"/ready")
	assertStatusNotFound(t, baseURL+"/orders/summary")
	assertStatusNotFound(t, baseURL+"/ping")
}

func TestStartREST_StandardGroupE2E(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)

	errCh := startRESTAsync(newStandardE2EBlueprint())
	baseURL := fmt.Sprintf("http://%s:%d", testHost, port)

	assertStandardRoutes(t, baseURL)
	assertStandardProbesAndMissingPaths(t, baseURL)

	signalAndWait(t, errCh)
}
