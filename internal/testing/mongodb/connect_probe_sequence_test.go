package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
)

const (
	errProbePrefix = "failed to probe mongodb"
	errStep2Substr = "failed to list collections without maxTimeMS"
	errStep3Substr = "failed to list collections with maxTimeMS"
	errStep4Substr = "failed to probe document in collection"
)

func setupMockProbeSequence(expectedErr error) func() {
	prevPing := SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	})
	prevProbe := SetMockProbe(func(context.Context, *client.Client, string) error {
		return expectedErr
	})
	return func() {
		prevProbe()
		prevPing()
	}
}

func assertProbeError(t *testing.T, err error, substr string) {
	t.Helper()
	require.Error(t, err)
	assert.Contains(t, err.Error(), errProbePrefix)
	assert.Contains(t, err.Error(), substr)
}

func TestConnect_ProbeSequence_Step2_Failure(t *testing.T) {
	cleanup := setupMockProbeSequence(errors.New(errStep2Substr + ": deadline"))
	defer cleanup()

	tc, err := Connect(context.Background())
	assert.Nil(t, tc)
	assertProbeError(t, err, errStep2Substr)
}

func TestConnect_ProbeSequence_Step3_Failure(t *testing.T) {
	cleanup := setupMockProbeSequence(errors.New(errStep3Substr + ": server error"))
	defer cleanup()

	tc, err := Connect(context.Background())
	assert.Nil(t, tc)
	assertProbeError(t, err, errStep3Substr)
}

func TestConnect_ProbeSequence_Step4_Failure(t *testing.T) {
	cleanup := setupMockProbeSequence(errors.New(errStep4Substr + " \"users\": err"))
	defer cleanup()

	tc, err := Connect(context.Background())
	assert.Nil(t, tc)
	assertProbeError(t, err, errStep4Substr)
}

func TestConnect_ProbeSequence_Success(t *testing.T) {
	cleanup := setupMockProbeSequence(nil)
	defer cleanup()

	tc, err := Connect(context.Background())
	require.NoError(t, err)
	require.NotNil(t, tc)
	assert.NoError(t, tc.Close())
}

func TestConnectURI_ProbeSequence_Failure(t *testing.T) {
	cleanup := setupMockProbeSequence(errors.New(errStep2Substr + ": context canceled"))
	defer cleanup()

	tc, err := ConnectURI(context.Background(), "mongodb://127.0.0.1:27017/testdb")
	assert.Nil(t, tc)
	assertProbeError(t, err, errStep2Substr)
}
