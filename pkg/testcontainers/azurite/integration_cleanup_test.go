//go:build testcontainers

package azurite_test

import (
	"context"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
	"github.com/nawaphonOHM/whatever/v3/pkg/testcontainers/azurite"
)

const (
	cleanupTimeout = 30 * time.Second
)

func registerEmulatorCleanup(t *testing.T, emulator *azurite.Container) func() {
	t.Helper()
	alreadyTerminated := false
	if emulator != nil {
		t.Cleanup(func() {
			cleanupEmulator(t, emulator, &alreadyTerminated)
		})
	}
	return func() {
		alreadyTerminated = true
	}
}

func cleanupEmulator(t *testing.T, emulator *azurite.Container, alreadyTerminated *bool) {
	t.Helper()
	if *alreadyTerminated {
		return
	}
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cleanupCancel()
	if err := emulator.Terminate(cleanupCtx); err != nil {
		t.Errorf("terminate Azurite container: %v", err)
	}
}

func registerS3ObjectCleanup(t *testing.T, client *storage.Client, bucket, key string) {
	t.Helper()
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), cleanupTimeout)
		defer cleanupCancel()
		_, deleteErr := client.DeleteObject(cleanupCtx, storage.DeleteObjectInput{Bucket: bucket, Key: key})
		if deleteErr != nil {
			t.Errorf("delete S3 gateway test object during cleanup: %v", deleteErr)
		}
	})
}
