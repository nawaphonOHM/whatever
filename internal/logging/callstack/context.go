package callstack

import (
	"context"
	"sync"
)

type stackContextKey struct{}

var callStackKey = stackContextKey{}

// WithCallStack returns a new context with the given CallStack attached.
func WithCallStack(ctx context.Context, stack *CallStack) context.Context {
	targetCtx := ctx
	if targetCtx == nil {
		targetCtx = context.Background()
	}
	return context.WithValue(targetCtx, callStackKey, stack)
}

// FromContext extracts the CallStack attached to context, or nil if none.
func FromContext(ctx context.Context) *CallStack {
	if ctx == nil {
		return nil
	}
	if cs, ok := ctx.Value(callStackKey).(*CallStack); ok {
		return cs
	}
	return nil
}

// EnsureContext returns the existing CallStack or creates and attaches a new CallStack.
func EnsureContext(ctx context.Context) (context.Context, *CallStack) {
	targetCtx := ctx
	if targetCtx == nil {
		targetCtx = context.Background()
	}
	if cs := FromContext(targetCtx); cs != nil {
		return targetCtx, cs
	}
	cs := New()
	return WithCallStack(targetCtx, cs), cs
}

var (
	defaultStackMu sync.RWMutex
	defaultStack   = New()
)

// DefaultStack returns the package-level default CallStack.
func DefaultStack() *CallStack {
	defaultStackMu.RLock()
	defer defaultStackMu.RUnlock()
	return defaultStack
}

// SetDefaultStack sets the package-level default CallStack.
func SetDefaultStack(cs *CallStack) {
	defaultStackMu.Lock()
	defer defaultStackMu.Unlock()
	if cs != nil {
		defaultStack = cs
	} else {
		defaultStack = New()
	}
}
