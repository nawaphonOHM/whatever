package callstack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallStack_ContextHelpers(t *testing.T) {
	assertEmptyContexts(t)
	assertAttachedContext(t)
	assertCreatedContext(t)
}

func assertEmptyContexts(t *testing.T) {
	t.Helper()
	assert.Nil(t, FromContext(context.TODO()))
	assert.Nil(t, FromContext(context.Background()))
}

func assertAttachedContext(t *testing.T) {
	t.Helper()
	cs := New()
	ctx := WithCallStack(context.Background(), cs)
	assert.Equal(t, cs, FromContext(ctx))
	ctx2, cs2 := EnsureContext(ctx)
	assert.Equal(t, cs, cs2)
	assert.Equal(t, ctx, ctx2)
}

func assertCreatedContext(t *testing.T) {
	t.Helper()
	ctx, cs := EnsureContext(context.Background())
	require.NotNil(t, cs)
	assert.Equal(t, cs, FromContext(ctx))
}

func TestCallStack_NilReceiverSafety(t *testing.T) {
	var cs *CallStack
	cs.Push(&Frame{Function: "noop"})
	assert.Equal(t, 0, cs.Len())
	_, ok := cs.Pop()
	assert.False(t, ok)
	assert.Nil(t, cs.Frames())
	assert.Nil(t, cs.Clone())
}
