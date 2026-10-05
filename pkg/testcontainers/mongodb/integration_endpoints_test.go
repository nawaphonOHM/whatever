//go:build testcontainers

package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/mongodb"
)

func assertContainerGetters(t *testing.T, c *mongodb.Container) {
	t.Helper()
	assert.Equal(t, testDBName, c.Database())
	assert.Equal(t, testUser, c.Username())
	assert.Equal(t, testPass, c.Password())
	assert.NotNil(t, c.RawContainer())
}

func assertHostAndPort(ctx context.Context, t *testing.T, c *mongodb.Container) (string, int) {
	t.Helper()
	host, hErr := c.Host(ctx)
	port, pErr := c.Port(ctx)
	require.NoError(t, hErr)
	require.NoError(t, pErr)
	assert.NotEmpty(t, host)
	assert.Greater(t, port, 0)
	return host, port
}

func assertConnString(ctx context.Context, t *testing.T, c *mongodb.Container) {
	t.Helper()
	connStr, err := c.ConnectionString(ctx)
	require.NoError(t, err)
	assert.Contains(t, connStr, "mongodb://")
}

func assertContainerConfig(ctx context.Context, t *testing.T, c *mongodb.Container) {
	t.Helper()
	cfg, err := c.Config(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.Host)
	assert.Greater(t, cfg.Port, 0)
	assert.Equal(t, testDBName, cfg.Database)
	assert.Equal(t, testUser, cfg.Username)
	assert.Equal(t, testPass, cfg.Password)
}

func TestContainer_Integration_Endpoints(t *testing.T) {
	requireSharedReady(t)
	ctx := context.Background()
	assertContainerGetters(t, sharedContainer)
	assertHostAndPort(ctx, t, sharedContainer)
	assertConnString(ctx, t, sharedContainer)
	assertContainerConfig(ctx, t, sharedContainer)
}
