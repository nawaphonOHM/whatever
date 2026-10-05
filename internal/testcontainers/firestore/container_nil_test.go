package firestore

import (
	"context"
	"errors"
	"testing"
)

func assertNilError(t *testing.T, name string, err error) {
	t.Helper()
	if !errors.Is(err, ErrNilContainer) {
		t.Errorf("%s expected ErrNilContainer, got %v", name, err)
	}
}

func assertNotRunningError(t *testing.T, name string, err error) {
	t.Helper()
	if !errors.Is(err, ErrContainerNotRunning) {
		t.Errorf("%s expected ErrContainerNotRunning, got %v", name, err)
	}
}

func testNilEndpoints(ctx context.Context, t *testing.T, c *Container) {
	_, err := c.URI(ctx)
	assertNilError(t, "URI", err)
	_, err = c.Host(ctx)
	assertNilError(t, "Host", err)
	_, err = c.Port(ctx)
	assertNilError(t, "Port", err)
}

func testNilBridges(ctx context.Context, t *testing.T, c *Container) {
	assertNilError(t, "Terminate", c.Terminate(ctx))
}

func TestNilContainer_Methods(t *testing.T) {
	var c *Container
	ctx := context.Background()
	// False positive: these nil-receiver calls are intentional nil-safety coverage;
	// proof: TestNilContainer_Methods in container_nil_test.go.
	testNilEndpoints(ctx, t, c)
	testNilBridges(ctx, t, c)
}

func testNilGetters(t *testing.T, c *Container) {
	if c.ProjectID() != "" {
		t.Errorf("expected empty string getter for ProjectID, got %q", c.ProjectID())
	}
	if c.DatastoreMode() {
		t.Errorf("expected false getter for DatastoreMode, got %v", c.DatastoreMode())
	}
}

func TestNilContainer_Getters(t *testing.T) {
	var c *Container
	// False positive: these nil-receiver calls are intentional nil-safety coverage;
	// proof: TestNilContainer_Getters in container_nil_test.go.
	testNilGetters(t, c)
	if c.RawContainer() != nil {
		t.Error("expected nil RawContainer on nil container")
	}
}

func testNotRunningEndpoints(ctx context.Context, t *testing.T, c *Container) {
	_, err := c.URI(ctx)
	assertNotRunningError(t, "URI", err)
	_, err = c.Host(ctx)
	assertNotRunningError(t, "Host", err)
	_, err = c.Port(ctx)
	assertNotRunningError(t, "Port", err)
}

func testNotRunningBridges(ctx context.Context, t *testing.T, c *Container) {
	assertNotRunningError(t, "Terminate", c.Terminate(ctx))
}

func TestNotRunningContainer_Methods(t *testing.T) {
	c := &Container{}
	ctx := context.Background()
	testNotRunningEndpoints(ctx, t, c)
	testNotRunningBridges(ctx, t, c)
}
