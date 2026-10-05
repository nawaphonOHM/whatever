//go:build testcontainers

package firestore_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
)

const (
	httpProbeTimeout = 5 * time.Second
)

func assertSharedEndpoints(ctx context.Context, t *testing.T, c *firestore.Container) {
	t.Helper()
	host, hErr := c.Host(ctx)
	port, pErr := c.Port(ctx)
	uri, uErr := c.URI(ctx)

	require.NoError(t, hErr)
	require.NoError(t, pErr)
	require.NoError(t, uErr)

	assert.NotEmpty(t, host)
	assert.Greater(t, port, 0)
	assert.Equal(t, fmt.Sprintf("%s:%d", host, port), uri)
}

func assertSharedMetadata(t *testing.T, c *firestore.Container) {
	t.Helper()
	assert.NotEmpty(t, c.ProjectID())
	assert.NotNil(t, c.RawContainer())
}

func assertTCPReachability(ctx context.Context, t *testing.T, c *firestore.Container) {
	t.Helper()
	uri, err := c.URI(ctx)
	require.NoError(t, err)

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", uri)
	require.NoError(t, err, "failed to establish TCP connection to Firestore emulator")
	require.NoError(t, conn.Close())
}

func queryEmulatorStatus(ctx context.Context, t *testing.T, uri string) int {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+uri+"/", nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: httpProbeTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to send HTTP request to Firestore emulator")
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	return resp.StatusCode
}

func assertHTTPReachability(ctx context.Context, t *testing.T, c *firestore.Container) {
	t.Helper()
	uri, err := c.URI(ctx)
	require.NoError(t, err)

	statusCode := queryEmulatorStatus(ctx, t, uri)
	assert.Equal(t, http.StatusOK, statusCode)
}

func assertFirestoreConfigIntegration(t *testing.T, c *firestore.Container) {
	t.Helper()
	firestoreTarget := fmt.Sprintf("%s.firestore.goog", c.ProjectID())
	assert.True(t, config.IsFirestoreURL(firestoreTarget), "IsFirestoreURL should detect firestore.goog host")

	firestoreURI := fmt.Sprintf("mongodb://user:pass@%s.firestore.goog:8080/testdb?authSource=admin", c.ProjectID())
	assert.True(t, config.IsFirestoreURL(firestoreURI), "IsFirestoreURL should detect firestore.goog URI")
}

func TestContainer_Integration_Ping(t *testing.T) {
	requireSharedReady(t)
	ctx := context.Background()

	assertSharedEndpoints(ctx, t, sharedContainer)
	assertSharedMetadata(t, sharedContainer)
	assertTCPReachability(ctx, t, sharedContainer)
	assertHTTPReachability(ctx, t, sharedContainer)
	assertFirestoreConfigIntegration(t, sharedContainer)
}
