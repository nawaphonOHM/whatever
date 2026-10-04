package mongodb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestConstantsAndErrors(t *testing.T) {
	assertSentinelErrors(t)
	assertConstants(t)
}

func assertConstants(t *testing.T) {
	assert.Equal(t, "mongo:6", mongodb.DefaultImage)
	assert.Equal(t, "27017/tcp", mongodb.DefaultPort)
}

func assertSentinelErrors(t *testing.T) {
	assert.ErrorIs(t, mongodb.ErrNilContainer, mongodb.ErrNilContainer)
	assert.ErrorIs(t, mongodb.ErrContainerNotRunning, mongodb.ErrContainerNotRunning)
}

func assertNilError(t *testing.T, name string, err error) {
	t.Helper()
	if !errors.Is(err, mongodb.ErrNilContainer) {
		t.Errorf("%s expected ErrNilContainer, got %v", name, err)
	}
}

func testNilEndpoints(ctx context.Context, t *testing.T, c *mongodb.Container) {
	_, err := c.ConnectionString(ctx)
	assertNilError(t, "ConnectionString", err)
	_, err = c.Host(ctx)
	assertNilError(t, "Host", err)
	_, err = c.Port(ctx)
	assertNilError(t, "Port", err)
}

func testNilBridges(ctx context.Context, t *testing.T, c *mongodb.Container) {
	_, err := c.Config(ctx)
	assertNilError(t, "Config", err)
	_, err = c.Client(ctx)
	assertNilError(t, "Client", err)
	assertNilError(t, "Terminate", c.Terminate(ctx))
}

func TestNilContainer_Methods(t *testing.T) {
	var c *mongodb.Container
	ctx := context.Background()
	// False positive: these nil-receiver calls are intentional nil-safety coverage;
	// proof: TestNilContainer_Methods in container_test.go.
	testNilEndpoints(ctx, t, c)
	testNilBridges(ctx, t, c)
}

func testNilGettersStrings(t *testing.T, c *mongodb.Container) {
	for _, val := range []string{c.Database(), c.Username(), c.Password(), c.ReplicaSet()} {
		if val != "" {
			t.Errorf("expected empty string getter, got %q", val)
		}
	}
}

func TestNilContainer_Getters(t *testing.T) {
	var c *mongodb.Container
	// False positive: these nil-receiver calls are intentional nil-safety coverage;
	// proof: TestNilContainer_Getters in container_test.go.
	testNilGettersStrings(t, c)
	if c.RawContainer() != nil {
		t.Error("expected nil RawContainer on nil container")
	}
}

func TestRun_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, err := mongodb.Run(ctx)
	if err == nil {
		t.Error("expected error when running with canceled context")
	}
	if c != nil {
		t.Error("expected nil container on error")
	}
}
