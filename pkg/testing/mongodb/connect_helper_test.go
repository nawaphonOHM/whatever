package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/pkg/testing/mongodb"
)

type mockTB struct {
	testing.TB
	fatalfMsg    string
	cleanups     []func()
	fatalfCalled bool
}

func (m *mockTB) Fatalf(format string, args ...any) {
	m.fatalfCalled = true
	m.fatalfMsg = fmt.Sprintf(format, args...)
}

func (m *mockTB) Cleanup(f func()) {
	m.cleanups = append(m.cleanups, f)
}

func (*mockTB) Helper() {}

func (*mockTB) Logf(string, ...any) {}

func setMockConnectionSuccess(t *testing.T) {
	t.Helper()
	t.Cleanup(mongodb.SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	}))
	t.Cleanup(mongodb.SetMockProbe(func(context.Context, *mongo.Client, string) error {
		return nil
	}))
}

func TestNewTestClient_NilTB(t *testing.T) {
	assert.PanicsWithError(t, mongodb.ErrNilTestingTB.Error(), func() {
		mongodb.NewTestClient(nil)
	})
	assert.PanicsWithError(t, mongodb.ErrNilTestingTB.Error(), func() {
		mongodb.NewTestClientURI(nil, "mongodb://localhost:27017")
	})
}

func TestNewTestClient_SuccessAndCleanup(t *testing.T) {
	setMockConnectionSuccess(t)
	mock := &mockTB{TB: t}
	client := mongodb.NewTestClient(mock, mongodb.WithDatabase("test_db"))
	require.NotNil(t, client)
	assert.False(t, mock.fatalfCalled)
	require.Len(t, mock.cleanups, 1)

	mock.cleanups[0]()
}

func TestNewTestClientURI_SuccessAndCleanup(t *testing.T) {
	setMockConnectionSuccess(t)
	mock := &mockTB{TB: t}
	client := mongodb.NewTestClientURI(mock, "mongodb://localhost:27017/uri_db")
	require.NotNil(t, client)
	assert.False(t, mock.fatalfCalled)
	require.Len(t, mock.cleanups, 1)

	mock.cleanups[0]()
}

func TestNewTestClient_Failure(t *testing.T) {
	mock := &mockTB{TB: t}
	client := mongodb.NewTestClient(mock, mongodb.WithPort(-1))
	assert.Nil(t, client)
	assert.True(t, mock.fatalfCalled)
	assert.Contains(t, mock.fatalfMsg, "failed to connect to mongodb")
}

func TestNewTestClientURI_Failure(t *testing.T) {
	mock := &mockTB{TB: t}
	client := mongodb.NewTestClientURI(mock, "")
	assert.Nil(t, client)
	assert.True(t, mock.fatalfCalled)
	assert.Contains(t, mock.fatalfMsg, "failed to connect to mongodb uri")
}
