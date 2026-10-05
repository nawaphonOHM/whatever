//go:build testcontainers

package firestore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
)

func runEphemeralContainer(ctx context.Context, t *testing.T) *firestore.Container {
	t.Helper()
	c, err := firestore.Run(ctx, firestore.WithProjectID("ephemeral-project"))
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

func assertContainerEndpoints(ctx context.Context, t *testing.T, c *firestore.Container) {
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

func assertContainerMetadata(t *testing.T, c *firestore.Container) {
	t.Helper()
	assert.NotEmpty(t, c.ProjectID())
	assert.False(t, c.DatastoreMode())
	assert.NotNil(t, c.RawContainer())
}

func assertTerminateCleanly(ctx context.Context, t *testing.T, c *firestore.Container) {
	t.Helper()
	require.NoError(t, c.Terminate(ctx))
}

func assertCanceledRunError(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, err := firestore.Run(ctx)
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestContainer_Integration_Lifecycle(t *testing.T) {
	ctx := context.Background()
	c := runEphemeralContainer(ctx, t)
	assertContainerEndpoints(ctx, t, c)
	assertContainerMetadata(t, c)
	assertTerminateCleanly(ctx, t, c)
	assertCanceledRunError(t)
}

func TestContainer_Integration_DatastoreModeLifecycle(t *testing.T) {
	ctx := context.Background()
	c, err := firestore.Run(ctx, firestore.WithDatastoreMode(true))
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.True(t, c.DatastoreMode())
	require.NoError(t, c.Terminate(ctx))
}
