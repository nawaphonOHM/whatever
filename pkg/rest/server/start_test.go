package server

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFreePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() {
		_ = l.Close()
	}()
	return l.Addr().(*net.TCPAddr).Port
}

func TestStartREST_Success(t *testing.T) {
	port := getFreePort(t)
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)
	t.Setenv("OHM9969_SERVER_HOST", "127.0.0.1")
	t.Setenv("OHM9969_SERVER_PORT", fmt.Sprintf("%d", port))
	t.Setenv("OHM9969_SERVER_SHUTDOWN_TIMEOUT", "2s")

	regs := []*RestAPIRegistration{
		{
			Version: 1,
			Prefix:  "/test",
			Apis: []*ExportableAPI{
				{
					Path:   "/ping",
					Method: GET,
					Handler: func(c *Context) response.Response {
						return response.OK("pong")
					},
				},
			},
		},
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- StartREST(regs)
	}()

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := fmt.Sprintf("http://127.0.0.1:%d/v1/test/ping", port)
	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err = client.Get(url)
		if err == nil {
			break
		}
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// Graceful shutdown
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))

	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server failed to shutdown in time")
	}
}

func TestStartREST_ReservedPathRejected(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)

	regs := []*RestAPIRegistration{
		{
			Version: 0,
			Prefix:  "",
			Apis: []*ExportableAPI{
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

	err := StartREST(regs)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrReservedPath)
}

func TestStartREST_DuplicateRejected(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)

	handler := func(c *Context) response.Response { return response.NoContent() }
	regs := []*RestAPIRegistration{
		{
			Version: 1,
			Prefix:  "/items",
			Apis: []*ExportableAPI{
				{Path: "", Method: GET, Handler: handler},
				{Path: "", Method: GET, Handler: handler},
			},
		},
	}

	err := StartREST(regs)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicateRoute)
}

func TestStartREST_NilRegistrationRejected(t *testing.T) {
	t.Setenv("OHM9969_GIN_MODE", gin.TestMode)

	err := StartREST([]*RestAPIRegistration{nil})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilRegistration)
}
