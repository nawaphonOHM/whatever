package callstack

import (
	"fmt"
	"sync"
	"time"
)

// Frame represents a single call frame in the call stack.
type Frame struct {
	Timestamp time.Time `json:"timestamp"`
	Function  string    `json:"function"`
	Package   string    `json:"package"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
}

// String returns a formatted single-line representation of a Frame.
func (f *Frame) String() string {
	if f == nil {
		return ""
	}
	pkg := f.Package
	if pkg != "" {
		pkg = pkg + "."
	}
	return fmt.Sprintf("%s%s (%s:%d) [%s]", pkg, f.Function, f.File, f.Line, f.Timestamp.Format(time.RFC3339Nano))
}

// CallStack provides thread-safe stack management.
type CallStack struct {
	frames []*Frame
	mu     sync.RWMutex
}

// New constructs an empty CallStack.
func New() *CallStack {
	return &CallStack{frames: make([]*Frame, 0)}
}

func ensureFrameTimestamp(frame *Frame) {
	if frame.Timestamp.IsZero() {
		frame.Timestamp = time.Now()
	}
}

// Push adds a frame to the top of the stack.
func (cs *CallStack) Push(frame *Frame) {
	if cs == nil || frame == nil {
		return
	}
	ensureFrameTimestamp(frame)
	cs.appendFrame(frame)
}

func (cs *CallStack) appendFrame(frame *Frame) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.frames = append(cs.frames, frame)
}

// Pop removes and returns the top frame of the stack.
func (cs *CallStack) Pop() (*Frame, bool) {
	if cs == nil {
		return nil, false
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if len(cs.frames) == 0 {
		return nil, false
	}
	lastIdx := len(cs.frames) - 1
	f := cs.frames[lastIdx]
	cs.frames = cs.frames[:lastIdx]
	return f, true
}

// Frames returns a copy of all frames currently on the stack.
func (cs *CallStack) Frames() []*Frame {
	if cs == nil {
		return nil
	}
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	copied := make([]*Frame, len(cs.frames))
	copy(copied, cs.frames)
	return copied
}

// Len returns the current frame count on the stack.
func (cs *CallStack) Len() int {
	if cs == nil {
		return 0
	}
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.frames)
}

// Clone creates a deep copy of the CallStack.
func (cs *CallStack) Clone() *CallStack {
	if cs == nil {
		return nil
	}
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	copied := make([]*Frame, len(cs.frames))
	copy(copied, cs.frames)
	return &CallStack{frames: copied}
}
