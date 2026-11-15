package azurite

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcazurite "github.com/testcontainers/testcontainers-go/modules/azure/azurite"
)

type fakeTestcontainer struct {
	testcontainers.Container
	terminateErr     error
	running          bool
	terminateInvoked bool
}

func (f *fakeTestcontainer) IsRunning() bool {
	return f.running
}

func (*fakeTestcontainer) PortEndpoint(ctx context.Context, port, protocol string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	mappedPorts := map[string]string{
		tcazurite.BlobPort:  "49100",
		tcazurite.QueuePort: "49101",
		tcazurite.TablePort: "49102",
	}
	mappedPort, ok := mappedPorts[port]
	if !ok {
		return "", fmt.Errorf("unexpected service port %q", port)
	}
	return protocol + "://127.0.0.1:" + mappedPort, nil
}

func (f *fakeTestcontainer) Terminate(context.Context, ...testcontainers.TerminateOption) error {
	f.terminateInvoked = true
	if f.terminateErr != nil {
		return f.terminateErr
	}
	f.running = false
	return nil
}

func newTestContainer(running bool) (*Container, *fakeTestcontainer) {
	raw := &fakeTestcontainer{running: running}
	return NewContainer(&tcazurite.Container{Container: raw}), raw
}

func TestContainerServiceURLs(t *testing.T) {
	ctx := context.Background()
	c, _ := newTestContainer(true)
	tests := []struct {
		name     string
		endpoint func(context.Context) (string, error)
		want     string
	}{
		{name: "blob", endpoint: c.BlobServiceURL, want: "http://127.0.0.1:49100"},
		{name: "queue", endpoint: c.QueueServiceURL, want: "http://127.0.0.1:49101"},
		{name: "table", endpoint: c.TableServiceURL, want: "http://127.0.0.1:49102"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.endpoint(ctx)
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestContainerServiceURLs_NilAndUnavailable(t *testing.T) {
	ctx := context.Background()
	var nilContainer *Container
	_, err := nilContainer.BlobServiceURL(ctx)
	require.ErrorIs(t, err, ErrNilContainer)

	c, _ := newTestContainer(false)
	_, err = c.BlobServiceURL(ctx)
	require.ErrorIs(t, err, ErrContainerUnavailable)

	assert.Nil(t, (*Container)(nil).RawContainer())
}

func TestContainerTerminate(t *testing.T) {
	c, raw := newTestContainer(true)
	require.NoError(t, c.Terminate(context.Background()))
	assert.True(t, raw.terminateInvoked)

	_, err := c.BlobServiceURL(context.Background())
	require.ErrorIs(t, err, ErrContainerUnavailable)
}

func TestContainerTerminate_ErrorKeepsContainerAvailable(t *testing.T) {
	terminateErr := errors.New("terminate failed")
	c, raw := newTestContainer(true)
	raw.terminateErr = terminateErr

	err := c.Terminate(context.Background())
	require.ErrorIs(t, err, terminateErr)
	_, err = c.BlobServiceURL(context.Background())
	require.NoError(t, err)
}
