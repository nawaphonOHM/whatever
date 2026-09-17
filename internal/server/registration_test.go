package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/health"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPMethod(t *testing.T) {
	tests := []struct {
		method   HTTPMethod
		expected string
		valid    bool
	}{
		{GET, "GET", true},
		{HEAD, "HEAD", true},
		{POST, "POST", true},
		{PUT, "PUT", true},
		{PATCH, "PATCH", true},
		{DELETE, "DELETE", true},
		{OPTIONS, "OPTIONS", true},
		{CONNECT, "CONNECT", true},
		{TRACE, "TRACE", true},
		{HTTPMethod(999), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.method.IsValid())
			if tt.valid {
				assert.Equal(t, tt.expected, tt.method.GetName())
				assert.Equal(t, tt.expected, tt.method.String())
			} else {
				assert.Panics(t, func() { tt.method.GetName() })
				assert.Contains(t, tt.method.String(), "UNKNOWN")
			}
		})
	}
}

func TestCalculateFullPath(t *testing.T) {
	tests := []struct {
		name     string
		version  ApiVersioning
		prefix   Pathz
		path     Pathz
		expected string
	}{
		{"empty all", 0, "", "", "/"},
		{"unversioned simple", 0, "/items", "/list", "/items/list"},
		{"unversioned missing slashes", 0, "items", "list", "/items/list"},
		{"unversioned trailing slash", 0, "/items/", "/list/", "/items/list"},
		{"versioned prefix", 1, "/users", "/:id", "/v1/users/:id"},
		{"explicit api in prefix is caller-owned", 1, "/api/users", "", "/v1/api/users"},
		{"versioned prefix already containing /v1", 1, "/v1/users", "", "/v1/users"},
		{"versioned prefix already containing /v2", 2, "/v2/items", "/list", "/v2/items/list"},
		{"versioned empty prefix with path", 1, "", "/items", "/v1/items"},
		{"reserved health exact match", 0, "", "/health", "/health"},
		{"reserved health in prefix", 0, "/health", "", "/health"},
		{"reserved ready exact match", 0, "", "/ready", "/ready"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateFullPath(tt.version, tt.prefix, tt.path)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestValidateRegistrations_Errors(t *testing.T) {
	dummyHandler := func(c *Context) response.Response { return response.NoContent() }

	t.Run("nil registration", func(t *testing.T) {
		regs := []*RestApiRegistration{nil}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilRegistration)
	})

	t.Run("nil api entry", func(t *testing.T) {
		regs := []*RestApiRegistration{
			{
				Version: 1,
				Prefix:  "/items",
				Apis:    []*ExportableApi{nil},
			},
		}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilAPI)
	})

	t.Run("nil handler", func(t *testing.T) {
		regs := []*RestApiRegistration{
			{
				Version: 1,
				Prefix:  "/items",
				Apis: []*ExportableApi{
					{
						Path:    "/test",
						Method:  GET,
						Handler: nil,
					},
				},
			},
		}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilHandler)
	})

	t.Run("invalid method", func(t *testing.T) {
		regs := []*RestApiRegistration{
			{
				Version: 1,
				Prefix:  "/items",
				Apis: []*ExportableApi{
					{
						Path:    "/test",
						Method:  HTTPMethod(999),
						Handler: dummyHandler,
					},
				},
			},
		}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidMethod)
	})

	t.Run("reserved /health collision in path", func(t *testing.T) {
		regs := []*RestApiRegistration{
			{
				Version: 0,
				Prefix:  "",
				Apis: []*ExportableApi{
					{
						Path:    "/health",
						Method:  GET,
						Handler: dummyHandler,
					},
				},
			},
		}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrReservedPath)
	})

	t.Run("reserved /ready collision in prefix", func(t *testing.T) {
		regs := []*RestApiRegistration{
			{
				Version: 0,
				Prefix:  "/ready",
				Apis: []*ExportableApi{
					{
						Path:    "",
						Method:  GET,
						Handler: dummyHandler,
					},
				},
			},
		}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrReservedPath)
	})

	t.Run("duplicate route registration", func(t *testing.T) {
		regs := []*RestApiRegistration{
			{
				Version: 1,
				Prefix:  "/items",
				Apis: []*ExportableApi{
					{
						Path:    "/list",
						Method:  GET,
						Handler: dummyHandler,
					},
					{
						Path:    "/list",
						Method:  GET,
						Handler: dummyHandler,
					},
				},
			},
		}
		_, err := ValidateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateRoute)
	})
}

func TestRegisterRoutesWithVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	regs := []*RestApiRegistration{
		{
			Version: 1,
			Prefix:  "/items",
			Apis: []*ExportableApi{
				{
					Path:   "",
					Method: GET,
					Handler: func(c *Context) response.Response {
						return response.OK([]string{"item1", "item2"})
					},
				},
			},
		},
	}

	err := RegisterRoutesWithVersion(engine, regs, "1.0.0")
	require.NoError(t, err)

	// Verify reserved /health
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var healthResp response.SuccessResponse[health.Status]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &healthResp))
	assert.True(t, healthResp.Success)
	assert.Equal(t, "up", healthResp.Data.Status)
	assert.Equal(t, "1.0.0", healthResp.Data.Version)

	// Verify reserved /ready
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/ready", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var readyResp response.SuccessResponse[health.Status]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &readyResp))
	assert.True(t, readyResp.Success)
	assert.Equal(t, "ready", readyResp.Data.Status)

	// Verify registered route /v1/items
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/v1/items", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var itemsResp response.SuccessResponse[[]string]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &itemsResp))
	assert.True(t, itemsResp.Success)
	assert.Equal(t, []string{"item1", "item2"}, itemsResp.Data)
}

func TestRegisterRoutes_MiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	var order []string

	mw1 := func(c *Context) {
		order = append(order, "mw1_before")
		c.Next()
		order = append(order, "mw1_after")
	}

	mw2 := func(c *Context) {
		order = append(order, "mw2_before")
		c.Next()
		order = append(order, "mw2_after")
	}

	regs := []*RestApiRegistration{
		{
			Version: 1,
			Prefix:  "/chain",
			Apis: []*ExportableApi{
				{
					Path:       "/test",
					Method:     GET,
					Middleware: []Middleware{mw1, mw2},
					Handler: func(c *Context) response.Response {
						order = append(order, "handler")
						return response.OK("done")
					},
				},
			},
		},
	}

	err := RegisterRoutes(engine, regs)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/chain/test", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, []string{"mw1_before", "mw2_before", "handler", "mw2_after", "mw1_after"}, order)
}
