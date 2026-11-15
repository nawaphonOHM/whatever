//go:build testcontainers

package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestLocalStackConnection(t *testing.T) {
	requireSharedLocalStack(t)
	require.NoError(t, sharedStorage.Ping(context.Background()))
	require.Equal(t, storage.ProviderS3Compatible, sharedStorage.Provider())
	require.Equal(t, localStackBucket, sharedStorage.DefaultBucket())
}
