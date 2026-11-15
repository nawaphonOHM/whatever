package azurite

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcazurite "github.com/testcontainers/testcontainers-go/modules/azure/azurite"
)

type runCapture struct {
	image       string
	customizers []testcontainers.ContainerCustomizer
}

type runOutcome struct {
	container  *Container
	err        error
	startupErr error
	fake       *fakeTestcontainer
}

type azuriteRunner func(
	context.Context,
	string,
	...testcontainers.ContainerCustomizer,
) (*tcazurite.Container, error)

func TestRunCleansUpPartiallyStartedContainer(t *testing.T) {
	startupErr := errors.New("startup failed")
	fake := &fakeTestcontainer{running: true}
	raw := &tcazurite.Container{Container: fake}
	state := runCapture{}
	runner := newTestRunner(raw, startupErr, &state)
	container, err := runWithRunner(
		context.Background(),
		runner,
		WithImage("custom-azurite:1"),
		WithEnv("KEY", "VALUE"),
	)
	assertRunFailedAndCleaned(t, runOutcome{
		container:  container,
		err:        err,
		startupErr: startupErr,
		fake:       fake,
	}, state)
}

func newTestRunner(raw *tcazurite.Container, startupErr error, state *runCapture) azuriteRunner {
	return func(
		ctx context.Context,
		image string,
		customizers ...testcontainers.ContainerCustomizer,
	) (*tcazurite.Container, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		state.image = image
		state.customizers = customizers
		return raw, startupErr
	}
}

func assertRunFailedAndCleaned(t *testing.T, outcome runOutcome, state runCapture) {
	t.Helper()
	require.Nil(t, outcome.container)
	require.ErrorIs(t, outcome.err, outcome.startupErr)
	assert.True(t, outcome.fake.terminateInvoked)
	assert.Equal(t, "custom-azurite:1", state.image)
	assert.Len(t, state.customizers, 1)
}
