package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupMockPingAndProbe(probeErr error) func() {
	cleanupPing := SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	})
	cleanupProbe := SetMockProbe(func(context.Context, *Client, string) error {
		return probeErr
	})
	return func() {
		cleanupProbe()
		cleanupPing()
	}
}

func TestSetMockProbe_OverridesAndRestores(t *testing.T) {
	customErr := errors.New("custom mock probe error")
	restore := SetMockProbe(func(context.Context, *Client, string) error {
		return customErr
	})
	defer restore()

	assert.Equal(t, customErr, probeClient(context.Background(), nil, "testdb"))

	restore()
	assert.Equal(t, ErrNilClient, probeClient(context.Background(), nil, "testdb"))
}

func TestSetMockProbe_NilRestoresDefault(t *testing.T) {
	restore := SetMockProbe(func(context.Context, *Client, string) error {
		return errors.New("temporary error")
	})
	defer restore()

	SetMockProbe(nil)
	assert.Equal(t, ErrNilClient, probeClient(context.Background(), nil, "testdb"))
}
