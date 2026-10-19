package mongodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewTestClient_NilTB(t *testing.T) {
	assert.PanicsWithError(t, ErrNilTestingTB.Error(), func() {
		NewTestClient(nil)
	})
	assert.PanicsWithError(t, ErrNilTestingTB.Error(), func() {
		NewTestClientURI(nil, "mongodb://localhost:27017")
	})
}

func TestNewTestClient_SuccessAndCleanup(t *testing.T) {
	setMockPingSuccess(t, nil)
	mock := &mockTB{TB: t}
	client := NewTestClient(mock, WithDatabase("test_db"))
	require.NotNil(t, client)
	assert.False(t, mock.fatalfCalled)
	require.Len(t, mock.cleanups, 1)

	mock.cleanups[0]()
}

func TestNewTestClient_Failure(t *testing.T) {
	mock := &mockTB{TB: t}
	client := NewTestClient(mock, WithPort(-1))
	assert.Nil(t, client)
	assert.True(t, mock.fatalfCalled)
	assert.Contains(t, mock.fatalfMsg, "failed to connect to mongodb")
}

func TestNewTestClientURI_Failure(t *testing.T) {
	mock := &mockTB{TB: t}
	client := NewTestClientURI(mock, "")
	assert.Nil(t, client)
	assert.True(t, mock.fatalfCalled)
	assert.Contains(t, mock.fatalfMsg, "failed to connect to mongodb uri")
}

func TestConnectURI_EmptyURI(t *testing.T) {
	ctx := context.Background()
	client, err := ConnectURI(ctx, "")
	require.ErrorIs(t, err, ErrEmptyURI)
	assert.Nil(t, client)

	client, err = ConnectURI(ctx, "   ")
	require.ErrorIs(t, err, ErrEmptyURI)
	assert.Nil(t, client)
}

func TestConnectURI_ExtractDatabase(t *testing.T) {
	setMockPingSuccess(t, nil)
	ctx := context.Background()
	client, err := ConnectURI(ctx, "mongodb://localhost:27017/extracted_db")
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "extracted_db", client.Database().Name())
}

func TestConnectURI_OverrideDatabase(t *testing.T) {
	setMockPingSuccess(t, nil)
	ctx := context.Background()
	client, err := ConnectURI(
		ctx,
		"mongodb://localhost:27017/extracted_db",
		WithDatabase("overridden_db"),
	)
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "overridden_db", client.Database().Name())
}
