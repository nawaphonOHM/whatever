//go:build testcontainers

package firestore_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
)

const (
	testProjectID = "test-firestore-project"
)

var (
	sharedContainer *firestore.Container
	initErr         error
)

func setupSharedContainer() error {
	ctx := context.Background()
	var err error
	sharedContainer, err = firestore.Run(ctx, firestore.WithProjectID(testProjectID))
	return err
}

func teardownSharedContainer() {
	if sharedContainer == nil {
		return
	}
	if err := sharedContainer.Terminate(context.Background()); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
}

func TestMain(m *testing.M) {
	initErr = setupSharedContainer()
	defer teardownSharedContainer()
	m.Run()
}

func requireSharedReady(t *testing.T) {
	t.Helper()
	require.NoError(t, initErr, "shared container failed to initialize")
	require.NotNil(t, sharedContainer, "shared container is nil")
}
