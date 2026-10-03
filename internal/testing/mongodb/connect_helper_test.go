package mongodb

import (
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
	mock := &mockTB{TB: t}
	client := NewTestClient(mock, WithPing(false), WithDatabase("test_db"))
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
