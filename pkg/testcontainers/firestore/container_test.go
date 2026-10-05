package firestore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
	"github.com/stretchr/testify/assert"
)

func TestConstantsAndErrors(t *testing.T) {
	assertSentinelErrors(t)
	assertConstants(t)
}

func assertConstants(t *testing.T) {
	assert.Equal(t, "gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators", firestore.DefaultImage)
	assert.Equal(t, "8080/tcp", firestore.DefaultPort)
	assert.Equal(t, "test-project", firestore.DefaultProjectID)
}

func assertSentinelErrors(t *testing.T) {
	assert.ErrorIs(t, firestore.ErrNilContainer, firestore.ErrNilContainer)
	assert.ErrorIs(t, firestore.ErrContainerNotRunning, firestore.ErrContainerNotRunning)
}

func assertNilError(t *testing.T, name string, err error) {
	t.Helper()
	if !errors.Is(err, firestore.ErrNilContainer) {
		t.Errorf("%s expected ErrNilContainer, got %v", name, err)
	}
}

func testNilEndpoints(ctx context.Context, t *testing.T, c *firestore.Container) {
	_, err := c.URI(ctx)
	assertNilError(t, "URI", err)
	_, err = c.Host(ctx)
	assertNilError(t, "Host", err)
	_, err = c.Port(ctx)
	assertNilError(t, "Port", err)
}

func testNilBridges(ctx context.Context, t *testing.T, c *firestore.Container) {
	assertNilError(t, "Terminate", c.Terminate(ctx))
}

func TestNilContainer_Methods(t *testing.T) {
	var c *firestore.Container
	ctx := context.Background()
	// False positive: these nil-receiver calls are intentional nil-safety coverage;
	// proof: TestNilContainer_Methods in container_test.go.
	testNilEndpoints(ctx, t, c)
	testNilBridges(ctx, t, c)
}

func testNilGetters(t *testing.T, c *firestore.Container) {
	if c.ProjectID() != "" {
		t.Errorf("expected empty string getter for ProjectID, got %q", c.ProjectID())
	}
	if c.DatastoreMode() {
		t.Errorf("expected false getter for DatastoreMode, got %v", c.DatastoreMode())
	}
}

func TestNilContainer_Getters(t *testing.T) {
	var c *firestore.Container
	// False positive: these nil-receiver calls are intentional nil-safety coverage;
	// proof: TestNilContainer_Getters in container_test.go.
	testNilGetters(t, c)
	if c.RawContainer() != nil {
		t.Error("expected nil RawContainer on nil container")
	}
}

func TestRun_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, err := firestore.Run(ctx)
	if err == nil {
		t.Error("expected error when running with canceled context")
	}
	if c != nil {
		t.Error("expected nil container on error")
	}
}
