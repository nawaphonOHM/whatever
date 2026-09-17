package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	intcfg "github.com/nawaphonOHM/whatever/internal/config"
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

func TestNew(t *testing.T) {
	cfg := Config{
		Host:            "127.0.0.1",
		Port:            8888,
		Mode:            gin.TestMode,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		IdleTimeout:     30 * time.Second,
		ShutdownTimeout: 3 * time.Second,
	}

	srv := New(&cfg)
	require.NotNil(t, srv)
	assert.Equal(t, &cfg, srv.Config)
	assert.NotNil(t, srv.Engine)
	assert.NotNil(t, srv.httpServer)
	assert.Equal(t, "127.0.0.1:8888", srv.httpServer.Addr)
}

func TestServer_StartAndGracefulShutdown(t *testing.T) {
	port := getFreePort(t)
	cfg := Config{
		Host:            "127.0.0.1",
		Port:            port,
		Mode:            gin.TestMode,
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		ShutdownTimeout: 2 * time.Second,
	}

	srv := New(&cfg)
	srv.Engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)

	go func() {
		errChan <- srv.Start(ctx)
	}()

	// Wait for server to accept connections
	url := fmt.Sprintf("http://127.0.0.1:%d/ping", port)
	var resp *http.Response
	var err error
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err = http.Get(url)
		if err == nil {
			break
		}
	}
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "pong", string(body))

	// Cancel context to trigger graceful shutdown
	cancel()

	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server failed to shutdown within timeout")
	}
}

func TestServer_PortInUseError(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() {
		_ = l.Close()
	}()

	port := l.Addr().(*net.TCPAddr).Port

	cfg := Config{
		Host:            "127.0.0.1",
		Port:            port,
		Mode:            gin.TestMode,
		ShutdownTimeout: 1 * time.Second,
	}

	srv := New(&cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = srv.Start(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "http server failed to start")
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, gin.ReleaseMode, cfg.Mode)
	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 60*time.Second, cfg.IdleTimeout)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("OHM9969_SERVER_HOST", "127.0.0.1")
	t.Setenv("OHM9969_SERVER_PORT", "9090")
	t.Setenv("OHM9969_GIN_MODE", "test")
	t.Setenv("OHM9969_APP_VERSION", "v1.2.3")
	t.Setenv("OHM9969_SERVER_READ_TIMEOUT", "15s")
	t.Setenv("OHM9969_SERVER_WRITE_TIMEOUT", "20s")
	t.Setenv("OHM9969_SERVER_IDLE_TIMEOUT", "45s")
	t.Setenv("OHM9969_SERVER_SHUTDOWN_TIMEOUT", "5s")

	cfg, err := intcfg.Load[Config]()
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "test", cfg.Mode)
	assert.Equal(t, "v1.2.3", cfg.AppVersion)
	assert.Equal(t, 15*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 20*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 45*time.Second, cfg.IdleTimeout)
	assert.Equal(t, 5*time.Second, cfg.ShutdownTimeout)
}

func TestServer_SetupDefaultMiddlewares_NoRoute(t *testing.T) {
	srv := New(DefaultConfig())
	srv.SetupDefaultMiddlewares()

	req := httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
	w := httptest.NewRecorder()
	srv.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}

func TestServer_SetupDefaultMiddlewares_NoMethod(t *testing.T) {
	srv := New(DefaultConfig())
	srv.SetupDefaultMiddlewares()

	srv.Engine.GET("/test-endpoint", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/test-endpoint", nil)
	w := httptest.NewRecorder()
	srv.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "METHOD_NOT_ALLOWED")
}

func TestServer_ShutdownDirect(t *testing.T) {
	port := getFreePort(t)
	cfg := Config{
		Host:            "127.0.0.1",
		Port:            port,
		Mode:            gin.TestMode,
		ShutdownTimeout: 2 * time.Second,
	}

	srv := New(&cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start(ctx)
	}()

	time.Sleep(50 * time.Millisecond)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()

	err := srv.Shutdown(shutdownCtx)
	assert.NoError(t, err)
}
