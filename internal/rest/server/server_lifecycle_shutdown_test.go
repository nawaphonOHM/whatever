package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createSlowServer(t *testing.T, started, block chan struct{}) (*Server, int) {
	t.Helper()
	port := getFreePort(t)
	srv, err := New(&Config{
		Host:            testHost,
		Port:            port,
		Mode:            gin.TestMode,
		ShutdownTimeout: 20 * time.Millisecond,
	})
	require.NoError(t, err)

	srv.Engine.GET("/slow", func(c *gin.Context) {
		close(started)
		<-block
		c.String(http.StatusOK, "slow")
	})
	return srv, port
}

func triggerSlowRequest(t *testing.T, port int) {
	t.Helper()
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("http://%s:%d/slow", testHost, port),
		nil,
	)
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()
}

func waitForStarted(t *testing.T, started <-chan struct{}) {
	t.Helper()
	select {
	case <-started:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for slow request to start")
	}
}

type slowServerTest struct {
	cancel  context.CancelFunc
	errCh   <-chan error
	started chan struct{}
	block   chan struct{}
	port    int
}

func runSlowShutdown(t *testing.T, sst *slowServerTest) error {
	t.Helper()
	if sst == nil {
		return nil
	}
	go triggerSlowRequest(t, sst.port)
	waitForStarted(t, sst.started)

	sst.cancel()
	err := waitErr(t, sst.errCh)
	close(sst.block)
	return err
}

func assertShutdownTimeoutLog(t *testing.T, data []byte, err error) {
	t.Helper()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "graceful shutdown failed")

	entries := parseJSONLogs(data)
	failEntry := findLogEntry(entries, "REST server graceful shutdown failed")
	require.NotNil(t, failEntry)
	assert.Equal(t, "ERROR", failEntry["level"])
	assert.NotEmpty(t, failEntry["error"])
}

func TestServer_LifecycleLogging_GracefulShutdownTimeout(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	started, block := make(chan struct{}), make(chan struct{})
	srv, port := createSlowServer(t, started, block)

	cancel, errCh := startWithCancel(srv)
	sst := &slowServerTest{cancel: cancel, errCh: errCh, started: started, block: block, port: port}
	err := runSlowShutdown(t, sst)

	assertShutdownTimeoutLog(t, buf.Bytes(), err)
}
