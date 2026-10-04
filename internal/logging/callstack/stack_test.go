package callstack

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLine10 = 10
	testLine20 = 20
	testLine42 = 42
	testLine99 = 99
	testYear   = 2026
	testMonth  = 10
	testDay    = 4
)

func assertPushFrames(t *testing.T, cs *CallStack, f1, f2 *Frame) {
	cs.Push(f1)
	assert.Equal(t, 1, cs.Len())
	cs.Push(f2)
	assert.Equal(t, 2, cs.Len())
	frames := cs.Frames()
	require.Len(t, frames, 2)
	assert.Equal(t, "func1", frames[0].Function)
	assert.Equal(t, "func2", frames[1].Function)
	assert.False(t, frames[1].Timestamp.IsZero())
}

func assertPopFrames(t *testing.T, cs *CallStack) {
	popped, ok := cs.Pop()
	assert.True(t, ok)
	assert.Equal(t, "func2", popped.Function)
	assert.Equal(t, 1, cs.Len())

	popped, ok = cs.Pop()
	assert.True(t, ok)
	assert.Equal(t, "func1", popped.Function)
	assert.Equal(t, 0, cs.Len())

	_, ok = cs.Pop()
	assert.False(t, ok)
}

func TestCallStack_PushPop(t *testing.T) {
	cs := New()
	assert.Equal(t, 0, cs.Len())
	_, ok := cs.Pop()
	assert.False(t, ok)

	f1 := &Frame{Function: "func1", Package: "pkg1", File: "file1.go", Line: testLine10, Timestamp: time.Now()}
	f2 := &Frame{Function: "func2", Package: "pkg2", File: "file2.go", Line: testLine20}

	assertPushFrames(t, cs, f1, f2)
	assertPopFrames(t, cs)
}

func TestCallStack_Clone(t *testing.T) {
	cs := New()
	cs.Push(&Frame{Function: "orig", Package: "pkg", File: "orig.go", Line: 1})

	cloned := cs.Clone()
	require.NotNil(t, cloned)
	assert.Equal(t, 1, cloned.Len())

	cloned.Push(&Frame{Function: "cloned", Package: "pkg", File: "cloned.go", Line: 2})
	assert.Equal(t, 2, cloned.Len())
	assert.Equal(t, 1, cs.Len())
}

func TestCallStack_String(t *testing.T) {
	var nilStack *CallStack
	assert.Equal(t, "<nil stack>", nilStack.String())

	emptyStack := New()
	assert.Equal(t, "<empty stack>", emptyStack.String())

	cs := New()
	cs.Push(&Frame{
		Function:  "TestFunc",
		Package:   "callstack",
		File:      "/test/path.go",
		Line:      testLine42,
		Timestamp: time.Date(testYear, testMonth, testDay, 2, 0, 0, 0, time.UTC),
	})

	str := cs.String()
	assert.Contains(t, str, "CallStack:")
	assert.Contains(t, str, "callstack.TestFunc")
	assert.Contains(t, str, "/test/path.go:42")
}

func TestFrame_String(t *testing.T) {
	var nilFrame *Frame
	assert.Equal(t, "", nilFrame.String())

	f := &Frame{
		Function:  "Execute",
		Package:   "service",
		File:      "exec.go",
		Line:      testLine99,
		Timestamp: time.Date(testYear, testMonth, testDay, 2, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, "service.Execute (exec.go:99) [2026-10-04T02:00:00Z]", f.String())
}
