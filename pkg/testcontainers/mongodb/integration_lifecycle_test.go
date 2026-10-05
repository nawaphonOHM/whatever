//go:build testcontainers

package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/mongodb"
)

func runEphemeralContainer(ctx context.Context, t *testing.T) *mongodb.Container {
	t.Helper()
	c, err := mongodb.Run(ctx, mongodb.WithDatabase("ephemeral_db"))
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

func assertEphemeralHealthy(ctx context.Context, t *testing.T, c *mongodb.Container) {
	t.Helper()
	client, err := c.Client(ctx)
	require.NoError(t, err)
	require.NoError(t, client.Ping(ctx))
}

func assertTerminateCleanly(ctx context.Context, t *testing.T, c *mongodb.Container) {
	t.Helper()
	require.NoError(t, c.Terminate(ctx))
}

func assertCanceledRunError(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, err := mongodb.Run(ctx)
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestContainer_Integration_Lifecycle(t *testing.T) {
	ctx := context.Background()
	c := runEphemeralContainer(ctx, t)
	assertEphemeralHealthy(ctx, t, c)
	assertTerminateCleanly(ctx, t, c)
	assertCanceledRunError(t)
}
