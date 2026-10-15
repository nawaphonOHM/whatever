package client_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	testConnectTimeout = 300 * time.Millisecond
	maxTimeoutWait     = 5 * time.Second
	pendingConnBuffer  = 8
)

func acceptAll(ln net.Listener, conns chan<- net.Conn) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		conns <- conn
	}
}

func closeAll(t *testing.T, ln net.Listener, conns chan net.Conn) {
	t.Helper()
	assert.NoError(t, ln.Close())
	for len(conns) > 0 {
		assert.NoError(t, (<-conns).Close())
	}
}

// silentTLSServer accepts TCP connections and never answers the TLS handshake.
func silentTLSServer(t *testing.T) string {
	t.Helper()
	lc := net.ListenConfig{}
	ln, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	conns := make(chan net.Conn, pendingConnBuffer)
	go acceptAll(ln, conns)
	t.Cleanup(func() { closeAll(t, ln, conns) })
	return ln.Addr().String()
}

func TestConnect_ConnectTimeoutBoundsTLSHandshake(t *testing.T) {
	t.Parallel()

	cfg := endpointConfig("https://" + silentTLSServer(t))
	cfg.ConnectTimeout = testConnectTimeout
	c, err := client.ConnectWithConfig(t.Context(), cfg,
		client.WithSkipProbe(true), client.WithDriverOptions(disableRetries))
	require.NoError(t, err)

	start := time.Now()
	_, err = c.HeadObject(t.Context(), contracts.HeadObjectInput{Key: testObjectKey})
	require.Error(t, err)
	assert.Less(t, time.Since(start), maxTimeoutWait)
}

func connectBlocked(t *testing.T) *client.Client {
	t.Helper()
	cfg := endpointConfig(testHTTPSEndpoint)
	cfg.RequestTimeout = shortRequestTimeout
	c, err := client.ConnectWithConfig(t.Context(), cfg, connectOptionsFn(blockUntilDone)...)
	require.NoError(t, err)
	return c
}

func TestConnect_RequestTimeoutBoundsRequests(t *testing.T) {
	t.Parallel()

	c := connectBlocked(t)
	start := time.Now()
	_, headErr := c.HeadObject(t.Context(), contracts.HeadObjectInput{Key: testObjectKey})
	_, putErr := c.PutObject(t.Context(), contracts.PutObjectInput{Key: testObjectKey})
	assert.Error(t, headErr)
	assert.Error(t, putErr)
	assert.Less(t, time.Since(start), maxTimeoutWait)
}
