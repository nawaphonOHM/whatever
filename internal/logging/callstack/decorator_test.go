package callstack

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecorate(t *testing.T) {
	SetDefaultStack(New())
	assert.Equal(t, 0, DefaultStack().Len())

	var innerLen int
	var innerFrame *Frame
	decorated := Decorate("testTarget", func() {
		innerLen = DefaultStack().Len()
		frames := DefaultStack().Frames()
		if len(frames) > 0 {
			innerFrame = frames[len(frames)-1]
		}
	})

	decorated()

	assert.Equal(t, 1, innerLen)
	assert.Equal(t, "testTarget", innerFrame.Function)
	assert.Contains(t, innerFrame.Package, "callstack")
	assert.Equal(t, 0, DefaultStack().Len())
}

func TestDecorate_Nested(t *testing.T) {
	SetDefaultStack(New())

	var nestedFrames []*Frame

	inner := Decorate("innerFunc", func() {
		nestedFrames = DefaultStack().Frames()
	})

	outer := Decorate("outerFunc", func() {
		inner()
	})

	outer()

	require.Len(t, nestedFrames, 2)
	assert.Equal(t, "outerFunc", nestedFrames[0].Function)
	assert.Equal(t, "innerFunc", nestedFrames[1].Function)
	assert.Equal(t, 0, DefaultStack().Len())
}

func TestDecorateErr(t *testing.T) {
	SetDefaultStack(New())
	expectedErr := errors.New("something failed")

	var innerLen int
	decorated := DecorateErr("failingFunc", func() error {
		innerLen = DefaultStack().Len()
		return expectedErr
	})

	err := decorated()
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 1, innerLen)
	assert.Equal(t, 0, DefaultStack().Len())
}

func TestDecorateContext(t *testing.T) {
	var capturedStack *CallStack
	var stackDepth int

	decorated := DecorateContext("contextFunc", func(ctx context.Context) error {
		capturedStack = FromContext(ctx)
		if capturedStack != nil {
			stackDepth = capturedStack.Len()
		}
		return nil
	})

	err := decorated(context.Background())
	assert.NoError(t, err)
	require.NotNil(t, capturedStack)
	assert.Equal(t, 1, stackDepth)
	assert.Equal(t, 0, capturedStack.Len())
}

func TestDecorateContext_Nested(t *testing.T) {
	var innerFrames []*Frame

	inner := DecorateContext("childFunc", func(ctx context.Context) error {
		cs := FromContext(ctx)
		require.NotNil(t, cs)
		innerFrames = cs.Frames()
		return nil
	})

	outer := DecorateContext("parentFunc", func(ctx context.Context) error {
		return inner(ctx)
	})

	err := outer(context.Background())
	assert.NoError(t, err)
	require.Len(t, innerFrames, 2)
	assert.Equal(t, "parentFunc", innerFrames[0].Function)
	assert.Equal(t, "childFunc", innerFrames[1].Function)
}
