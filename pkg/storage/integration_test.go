//go:build testcontainers

package storage_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/localstack"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

var (
	sharedLocalStack     *localstack.LocalStackContainer
	sharedLocalStackErr  error
	sharedStorage        *storage.Client
	sharedStorageCleanup func()
)

func TestMain(m *testing.M) {
	if os.Getenv(envStorageEagerHelper) != "1" {
		sharedLocalStackErr = setupSharedLocalStack()
	}
	m.Run()
	if sharedStorageCleanup != nil {
		sharedStorageCleanup()
	}
}

func requireSharedLocalStack(t *testing.T) {
	t.Helper()
	require.NoError(t, sharedLocalStackErr, "shared LocalStack fixture failed to initialize")
	require.NotNil(t, sharedLocalStack)
	require.NotNil(t, sharedStorage)
}
