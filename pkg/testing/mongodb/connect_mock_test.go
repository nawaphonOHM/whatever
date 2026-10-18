package mongodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	testOverrideDB  = "override_db"
	testTimeout50ms = 50 * time.Millisecond
)

func TestSetMockProbe_OverridesAndRestores(t *testing.T) {
	t.Cleanup(mongodb.SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	}))

	t.Run("override probe fails connection", func(t *testing.T) {
		customErr := errors.New("custom mock probe failure")
		restore := mongodb.SetMockProbe(func(ctx context.Context, c *mongodb.Client, db string) error {
			assert.NotNil(t, ctx)
			assert.NotNil(t, c)
			assert.Equal(t, testOverrideDB, db)
			return customErr
		})
		defer restore()

		client, err := mongodb.Connect(context.Background(), mongodb.WithDatabase(testOverrideDB))
		assert.Nil(t, client)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "custom mock probe failure")
	})

	t.Run("restore probe allows connection", func(t *testing.T) {
		restore := mongodb.SetMockProbe(func(context.Context, *mongodb.Client, string) error {
			return errors.New("temporary failure")
		})
		restore()

		successProbeCalled := false
		restoreSuccess := mongodb.SetMockProbe(func(ctx context.Context, c *mongodb.Client, db string) error {
			successProbeCalled = true
			assert.NotNil(t, ctx)
			assert.NotNil(t, c)
			assert.Equal(t, testOverrideDB, db)
			return nil
		})
		defer restoreSuccess()

		client, err := mongodb.Connect(context.Background(), mongodb.WithDatabase(testOverrideDB))
		require.NoError(t, err)
		assert.True(t, successProbeCalled)
		require.NoError(t, client.Close())
	})
}

func TestSetMockProbe_NilRestoresDefault(t *testing.T) {
	t.Cleanup(mongodb.SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	}))

	restore := mongodb.SetMockProbe(func(context.Context, *mongodb.Client, string) error {
		return errors.New("temporary probe error")
	})
	defer restore()

	// Setting nil should reset to defaultProbeClient
	mongodb.SetMockProbe(nil)

	// Since default probe requires live mongo server or valid collections, Connect to an unreachable port will fail
	ctx := context.Background()
	client, err := mongodb.Connect(ctx,
		mongodb.WithPort(testUnreachablePort),
		mongodb.WithServerSelectionTimeout(testTimeout50ms),
		mongodb.WithConnectTimeout(testTimeout50ms),
	)
	assert.Nil(t, client)
	require.Error(t, err)
}

func TestSetMockPing_OverridesAndRestores(t *testing.T) {
	t.Cleanup(mongodb.SetMockProbe(func(context.Context, *mongodb.Client, string) error {
		return nil
	}))

	customErr := errors.New("custom mock ping failure")
	restore := mongodb.SetMockPing(func(context.Context, *mongo.Client) error {
		return customErr
	})

	ctx := context.Background()
	client, err := mongodb.Connect(ctx, mongodb.WithDatabase("ping_db"))
	assert.Nil(t, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "custom mock ping failure")

	restore()

	// Setting nil should reset to defaultPingClient
	mongodb.SetMockPing(nil)
}
