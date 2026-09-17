package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/health"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromRegistrations_Success(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)
	t.Setenv("OHM9969_APP_VERSION", "2.1.0")
	t.Setenv("OHM9969_SERVER_PORT", fmt.Sprintf("%d", getFreePort(t)))

	regs := []*RestApiRegistration{
		{
			Version: 1,
			Prefix:  "/items",
			Apis: []*ExportableApi{
				{
					Path:   "",
					Method: GET,
					Handler: func(c *Context) response.Response {
						return response.OK([]string{"a", "b"})
					},
				},
			},
		},
	}

	srv, err := NewFromRegistrations(regs)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// Reserved /health carries AppVersion from env
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	srv.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var healthResp response.SuccessResponse[health.Status]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &healthResp))
	assert.True(t, healthResp.Success)
	assert.Equal(t, "up", healthResp.Data.Status)
	assert.Equal(t, "2.1.0", healthResp.Data.Version)

	// Consumer route is mounted
	req, _ = http.NewRequest(http.MethodGet, "/v1/items", nil)
	w = httptest.NewRecorder()
	srv.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewFromRegistrations_ReservedPathRejected(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)

	regs := []*RestApiRegistration{
		{
			Version: 0,
			Prefix:  "",
			Apis: []*ExportableApi{
				{
					Path:   "/health",
					Method: GET,
					Handler: func(c *Context) response.Response {
						return response.OK("override")
					},
				},
			},
		},
	}

	srv, err := NewFromRegistrations(regs)
	require.Error(t, err)
	assert.Nil(t, srv)
	assert.ErrorIs(t, err, ErrReservedPath)
}

func TestNewFromRegistrations_DuplicateRejected(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)

	handler := func(c *Context) response.Response { return response.NoContent() }
	regs := []*RestApiRegistration{
		{
			Version: 1,
			Prefix:  "/items",
			Apis: []*ExportableApi{
				{Path: "", Method: GET, Handler: handler},
				{Path: "", Method: GET, Handler: handler},
			},
		},
	}

	srv, err := NewFromRegistrations(regs)
	require.Error(t, err)
	assert.Nil(t, srv)
	assert.ErrorIs(t, err, ErrDuplicateRoute)
}

func TestStartREST_RunsAndShutsDown(t *testing.T) {
	port := getFreePort(t)
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)
	t.Setenv("OHM9969_SERVER_HOST", "127.0.0.1")
	t.Setenv("OHM9969_SERVER_PORT", fmt.Sprintf("%d", port))
	t.Setenv("OHM9969_SERVER_SHUTDOWN_TIMEOUT", "2s")

	regs := []*RestApiRegistration{
		{
			Version: 1,
			Prefix:  "/ping",
			Apis: []*ExportableApi{
				{
					Path:   "",
					Method: GET,
					Handler: func(c *Context) response.Response {
						return response.OK("pong")
					},
				},
			},
		},
	}

	srv, err := NewFromRegistrations(regs)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start(ctx)
	}()

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	// 1. Check health
	urlHealth := fmt.Sprintf("http://127.0.0.1:%d/health", port)
	var resp *http.Response
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err = client.Get(urlHealth)
		if err == nil {
			break
		}
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 2. Check mounted route
	urlPing := fmt.Sprintf("http://127.0.0.1:%d/v1/ping", port)
	resp, err = client.Get(urlPing)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 3. Cancel context and wait for clean shutdown
	cancel()
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server failed to stop in time")
	}
}

func TestStartREST_ValidationFailure(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)

	err := StartREST([]*RestApiRegistration{nil})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilRegistration)
}
