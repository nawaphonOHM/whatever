package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPMethod(t *testing.T) {
	cases := map[HTTPMethod]string{
		GET: http.MethodGet, HEAD: http.MethodHead,
		POST: http.MethodPost, PUT: http.MethodPut,
		PATCH: http.MethodPatch, DELETE: http.MethodDelete,
		OPTIONS: http.MethodOptions, CONNECT: http.MethodConnect,
		TRACE: http.MethodTrace,
	}
	for method, name := range cases {
		assert.Equal(t, name, method.GetName())
		assert.Equal(t, name, method.String())
		assert.True(t, method.IsValid())
	}
}

func TestHTTPMethod_Invalid(t *testing.T) {
	method := HTTPMethod(invalidMethod)
	assert.False(t, method.IsValid())
	assert.Contains(t, method.String(), "UNKNOWN")
	assert.Panics(t, func() { method.GetName() })
}

func TestRestAPIRegistrationStructure(t *testing.T) {
	api := &ExportableAPI{Path: "/status", Method: GET}
	reg := &RRestAPIRegistration{
		Prefix: "/api", Apis: []*ExportableAPI{api},
		Version: APIVersioning(1),
	}
	assert.Equal(t, Pathz("/api"), reg.Prefix)
	assert.Equal(t, APIVersioning(1), reg.Version)
	require.Len(t, reg.Apis, 1)
	assert.Equal(t, Pathz("/status"), reg.Apis[0].Path)
	assert.Equal(t, GET, reg.Apis[0].Method)
}

func TestHandlerAndMiddlewareContracts(t *testing.T) {
	called := false
	api := ExportableAPI{
		Handler:    func(Context) Response { return OK(testValue) },
		Middleware: []Middleware{func(Context) { called = true }},
	}
	api.Middleware[0](nil)
	response := api.Handler(nil)
	assert.True(t, called)
	assert.Equal(t, testStatus, response.StatusCode())
}
